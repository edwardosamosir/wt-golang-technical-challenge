package usecase

import (
	"context"
	"fmt"
	"golang-technical-challenge/internal/entity"
	"golang-technical-challenge/internal/model"
	"golang-technical-challenge/internal/model/converter"
	"golang-technical-challenge/internal/repository"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type InvoiceUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	InvoiceRepository *repository.InvoiceRepository
}

func NewInvoiceUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, invoiceRepository *repository.InvoiceRepository,
) *InvoiceUseCase {
	return &InvoiceUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		InvoiceRepository: invoiceRepository,
	}
}

func parseDateFromCell(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	dateFormats := []string{
		"2006-01-02",                                           // ISO
		"02/01/2006", "01/02/2006", "02-01-2006", "01-02-2006", // DMY and MDY with 4-digit year
		"02/01/06", "01/02/06", "02-01-06", "01-02-06", // DMY and MDY with 2-digit year
		"2/1/06", "1/2/06", "2-1-06", "1-2-06", // Single-digit day/month
	}

	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	if serial, err := strconv.ParseFloat(raw, 64); err == nil {
		base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		return base.AddDate(0, 0, int(serial)), nil
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", raw)
}

func (c *InvoiceUseCase) ImportInvoices(ctx context.Context, file *multipart.FileHeader) (any, error) {
	f, err := file.Open()
	if err != nil {
		c.Log.WithError(err).Error("Failed to open uploaded file")
		return nil, err
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		c.Log.WithError(err).Error("Failed to parse XLSX file")
		return nil, err
	}
	defer xlsx.Close()

	invoiceRows, err := xlsx.GetRows("invoice")
	if err != nil {
		c.Log.WithError(err).Error("Failed to read 'invoice' sheet")
		return nil, fmt.Errorf("cannot read invoice sheet: %w", err)
	}
	productRows, err := xlsx.GetRows("product sold")
	if err != nil {
		c.Log.WithError(err).Error("Failed to read 'product sold' sheet")
		return nil, fmt.Errorf("cannot read product sold sheet: %w", err)
	}

	errors := []model.ImportError{}
	invoiceMap := map[string]*entity.Invoice{}

	c.parseInvoiceRows(ctx, xlsx, invoiceRows, invoiceMap, &errors)
	c.parseProductRows(productRows, invoiceMap, &errors)

	tx := c.DB.WithContext(ctx).Begin()
	for _, invoice := range invoiceMap {
		if len(invoice.Products) == 0 {
			errors = append(errors, model.ImportError{InvoiceNo: invoice.InvoiceNo, Message: "No valid products for this invoice"})
			continue
		}
		if err := c.InvoiceRepository.Create(tx, invoice); err != nil {
			errors = append(errors, model.ImportError{InvoiceNo: invoice.InvoiceNo, Message: "Failed to save invoice"})
		}
	}
	tx.Commit()

	if len(errors) > 0 {
		c.Log.WithField("error_count", len(errors)).Warn("Import completed with errors")
		return errors, nil
	}

	invoiceNos := make([]string, 0, len(invoiceMap))
	for _, inv := range invoiceMap {
		invoiceNos = append(invoiceNos, inv.InvoiceNo)
	}

	invoices, err := c.InvoiceRepository.FindInvoicesByNumbers(c.DB.WithContext(ctx), invoiceNos)
	if err != nil {
		c.Log.WithError(err).Error("Failed to retrieve imported invoices")
		return nil, err
	}

	totalProfit := decimal.Zero
	totalCash := decimal.Zero
	for _, inv := range invoices {
		for _, p := range inv.Products {
			cost := p.TotalCost.Mul(decimal.NewFromInt(int64(p.Quantity)))
			price := p.TotalPrice.Mul(decimal.NewFromInt(int64(p.Quantity)))
			totalProfit = totalProfit.Add(price.Sub(cost))
			if inv.PaymentType == "CASH" {
				totalCash = totalCash.Add(price)
			}
		}
	}

	invoiceResponses := converter.InvoicesToResponseList(invoices)
	totalItems := int64(len(invoices))

	return &model.InvoiceListResponse{
		Invoices:    invoiceResponses,
		TotalProfit: totalProfit.StringFixed(2),
		TotalCash:   totalCash.StringFixed(2),
		Paging: model.PageMetadata{
			Page:      1,
			Size:      len(invoices),
			TotalItem: totalItems,
			TotalPage: 1,
		},
	}, nil
}

func (c *InvoiceUseCase) parseInvoiceRows(ctx context.Context, xlsx *excelize.File, rows [][]string, invoiceMap map[string]*entity.Invoice, errors *[]model.ImportError) {
	for i, row := range rows[1:] {
		rowNum := i + 2
		if len(row) < 5 {
			*errors = append(*errors, model.ImportError{InvoiceNo: fmt.Sprintf("row %d", rowNum), Message: "Missing invoice fields"})
			continue
		}

		invoiceNo := strings.TrimSpace(row[0])
		customer := row[2]
		sales := row[3]
		paymentType := strings.ToUpper(row[4])
		if paymentType != "CASH" && paymentType != "CREDIT" {
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Invalid payment type"})
			continue
		}

		notes := ""
		if len(row) > 5 {
			notes = row[5]
		}

		dateCell := fmt.Sprintf("B%d", rowNum)
		dateStr, _ := xlsx.GetCellValue("invoice", dateCell)
		parsedDate, err := parseDateFromCell(dateStr)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"row":       rowNum,
				"invoiceNo": invoiceNo,
				"date":      dateStr,
			}).Warn("Invalid date format")
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Invalid date format"})
			continue
		}

		if invoiceNo == "" || customer == "" || sales == "" || paymentType == "" {
			c.Log.WithFields(logrus.Fields{
				"row":       rowNum,
				"invoiceNo": invoiceNo,
			}).Warn("Required invoice fields missing")
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Required invoice fields are missing"})
			continue
		}

		existing := new(entity.Invoice)
		if err := c.InvoiceRepository.FindByInvoiceNo(c.DB.WithContext(ctx), existing, invoiceNo); err == nil {
			c.Log.WithField("invoice_no", invoiceNo).Warn("Duplicate invoice")
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Duplicate invoice"})
			continue
		}

		invoiceMap[invoiceNo] = &entity.Invoice{
			InvoiceNo:       invoiceNo,
			Date:            parsedDate,
			CustomerName:    customer,
			SalespersonName: sales,
			PaymentType:     paymentType,
			Notes:           &notes,
			Products:        []entity.Product{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
	}
}

func (c *InvoiceUseCase) parseProductRows(rows [][]string, invoiceMap map[string]*entity.Invoice, errors *[]model.ImportError) {
	for i, row := range rows[1:] {
		rowNum := i + 2
		if len(row) < 5 {
			*errors = append(*errors, model.ImportError{InvoiceNo: fmt.Sprintf("row %d", rowNum), Message: "Missing product fields"})
			continue
		}

		invoiceNo := strings.TrimSpace(row[0])
		item := row[1]
		qtyStr := row[2]
		costStr := row[3]
		priceStr := row[4]

		invoice, ok := invoiceMap[invoiceNo]
		if !ok {
			c.Log.WithFields(logrus.Fields{
				"row":       rowNum,
				"invoiceNo": invoiceNo,
			}).Warn("Product refers to unknown invoice")
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Product refers to unknown invoice"})
			continue
		}

		qty, err1 := strconv.Atoi(qtyStr)
		cost, err2 := decimal.NewFromString(costStr)
		price, err3 := decimal.NewFromString(priceStr)
		if err1 != nil || err2 != nil || err3 != nil {
			c.Log.WithFields(logrus.Fields{
				"row":       rowNum,
				"invoiceNo": invoiceNo,
				"qty":       qtyStr,
				"cost":      costStr,
				"price":     priceStr,
			}).Warn("Invalid product values")
			*errors = append(*errors, model.ImportError{InvoiceNo: invoiceNo, Message: "Invalid product values"})
			continue
		}

		invoice.Products = append(invoice.Products, entity.Product{
			InvoiceNo:  invoiceNo,
			ItemName:   item,
			Quantity:   qty,
			TotalCost:  cost,
			TotalPrice: price,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
}

func (c *InvoiceUseCase) GetInvoices(ctx context.Context, date string, page, size int) (*model.InvoiceListResponse, error) {
	if date == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "date parameter is required")
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset := (page - 1) * size

	tx := c.DB.WithContext(ctx)

	invoices, totalItems, err := c.InvoiceRepository.FindInvoicesByDate(tx, date, size, offset)
	if err != nil {
		c.Log.WithError(err).WithField("date", date).Error("Failed to fetch invoices")
		return nil, fiber.ErrInternalServerError
	}

	totalProfit, totalCash, err := c.InvoiceRepository.GetSummaryByDate(tx, date)
	if err != nil {
		c.Log.WithError(err).WithField("date", date).Error("Failed to calculate invoice summary")
		return nil, fiber.ErrInternalServerError
	}

	invoiceResponses := converter.InvoicesToResponseList(invoices)
	totalPages := (totalItems + int64(size) - 1) / int64(size)

	return &model.InvoiceListResponse{
		Invoices:    invoiceResponses,
		TotalProfit: totalProfit,
		TotalCash:   totalCash,
		Paging: model.PageMetadata{
			Page:      page,
			Size:      size,
			TotalItem: totalItems,
			TotalPage: totalPages,
		},
	}, nil
}

func (c *InvoiceUseCase) Create(ctx context.Context, request *model.CreateInvoiceRequest) (*model.InvoiceResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid create invoice payload")
		return nil, fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.Log.WithError(err).Warn("Invalid date format for create invoice")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
	}

	existing := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByInvoiceNo(tx, existing, request.InvoiceNo); err == nil {
		c.Log.WithField("invoice_no", request.InvoiceNo).Warn("Invoice already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Invoice already exists")
	} else if err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).WithField("invoice_no", request.InvoiceNo).Error("Failed to check existing invoice")
		return nil, fiber.ErrInternalServerError
	}

	invoice := &entity.Invoice{
		InvoiceNo:       request.InvoiceNo,
		Date:            date,
		CustomerName:    request.CustomerName,
		SalespersonName: request.SalespersonName,
		PaymentType:     request.PaymentType,
		Notes:           request.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	products := make([]entity.Product, 0, len(request.Products))
	for _, p := range request.Products {
		products = append(products, entity.Product{
			InvoiceNo:  request.InvoiceNo,
			ItemName:   p.ItemName,
			Quantity:   p.Quantity,
			TotalCost:  p.TotalCost,
			TotalPrice: p.TotalPrice,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
	invoice.Products = products

	if err := c.InvoiceRepository.Create(tx, invoice); err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoice.InvoiceNo).Error("Failed to create invoice")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoice.InvoiceNo).Error("Failed to commit invoice creation")
		return nil, fiber.ErrInternalServerError
	}

	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) Update(ctx context.Context, invoiceNo string, request *model.UpdateInvoiceRequest) (*model.InvoiceResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoiceNo).Warn("Invalid update invoice payload")
		return nil, fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByInvoiceNo(tx, invoice, invoiceNo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.WithField("invoice_no", invoiceNo).Warn("Invoice not found")
			return nil, fiber.ErrNotFound
		}
		c.Log.WithError(err).WithField("invoice_no", invoiceNo).Error("Failed to fetch invoice for update")
		return nil, fiber.ErrInternalServerError
	}

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoiceNo).Warn("Invalid date format for update invoice")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
	}

	invoice.Date = date
	invoice.CustomerName = request.CustomerName
	invoice.SalespersonName = request.SalespersonName
	invoice.PaymentType = request.PaymentType
	invoice.Notes = request.Notes
	invoice.UpdatedAt = time.Now()

	if err := tx.Where("invoice_no = ?", invoice.InvoiceNo).Delete(&entity.Product{}).Error; err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoice.InvoiceNo).Error("Failed to delete old products before update")
		return nil, fiber.ErrInternalServerError
	}

	newProducts := make([]entity.Product, 0, len(request.Products))
	for _, p := range request.Products {
		newProducts = append(newProducts, entity.Product{
			InvoiceNo:  invoice.InvoiceNo,
			ItemName:   p.ItemName,
			Quantity:   p.Quantity,
			TotalCost:  p.TotalCost,
			TotalPrice: p.TotalPrice,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
	invoice.Products = newProducts

	if err := tx.Save(invoice).Error; err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoice.InvoiceNo).Error("Failed to update invoice")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoice.InvoiceNo).Error("Failed to commit invoice update")
		return nil, fiber.ErrInternalServerError
	}

	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) Delete(ctx context.Context, request *model.DeleteInvoiceRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).WithField("invoice_no", request.InvoiceNo).Warn("Invalid delete invoice payload")
		return fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByInvoiceNo(tx, invoice, request.InvoiceNo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.WithField("invoice_no", request.InvoiceNo).Warn("Invoice not found")
			return fiber.ErrNotFound
		}
		c.Log.WithError(err).WithField("invoice_no", request.InvoiceNo).Warn("Invalid delete invoice payload")
		return fiber.ErrInternalServerError
	}

	if err := c.InvoiceRepository.Delete(tx, invoice); err != nil {
		c.Log.WithError(err).WithField("invoice_no", request.InvoiceNo).Error("Failed to delete invoice")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).WithField("invoice_no", request.InvoiceNo).Error("Failed to commit invoice deletion")
		return fiber.ErrInternalServerError
	}

	return nil
}
