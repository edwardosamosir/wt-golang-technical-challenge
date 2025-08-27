package usecase

import (
	"context"
	"golang-technical-challenge/internal/entity"
	"golang-technical-challenge/internal/model"
	"golang-technical-challenge/internal/model/converter"
	"golang-technical-challenge/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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
		c.Log.WithError(err).Error("Failed to fetch invoices")
		return nil, fiber.ErrInternalServerError
	}

	totalProfit, totalCash, err := c.InvoiceRepository.GetSummaryByDate(tx, date)
	if err != nil {
		c.Log.WithError(err).Error("Failed to calculate summary")
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
		c.Log.WithError(err).Error("Error validating request body")
		return nil, fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.Log.WithError(err).Error("Invalid date format")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
	}

	existing := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByInvoiceNo(tx, existing, request.InvoiceNo); err == nil {
		c.Log.WithField("invoice_no", request.InvoiceNo).Warn("Invoice already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Invoice already exists")
	} else if err != gorm.ErrRecordNotFound {
		c.Log.WithError(err).Error("Error checking existing invoice")
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
		c.Log.WithError(err).Error("Error creating invoice")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) Update(ctx context.Context, invoiceNo string, request *model.UpdateInvoiceRequest) (*model.InvoiceResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Error validating request body")
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
		c.Log.WithError(err).Error("Error fetching invoice")
		return nil, fiber.ErrInternalServerError
	}

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.Log.WithError(err).Error("Invalid date format")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
	}

	invoice.Date = date
	invoice.CustomerName = request.CustomerName
	invoice.SalespersonName = request.SalespersonName
	invoice.PaymentType = request.PaymentType
	invoice.Notes = request.Notes
	invoice.UpdatedAt = time.Now()

	if err := tx.Where("invoice_no = ?", invoice.InvoiceNo).Delete(&entity.Product{}).Error; err != nil {
		c.Log.WithError(err).Error("Error deleting old products")
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
		c.Log.WithError(err).Error("Error updating invoice")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) Delete(ctx context.Context, request *model.DeleteInvoiceRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Error validating request body")
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
		c.Log.WithError(err).Error("Error fetching invoice from DB")
		return fiber.ErrInternalServerError
	}

	if err := c.InvoiceRepository.Delete(tx, invoice); err != nil {
		c.Log.WithError(err).Error("Error deleting invoice")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Error deleting invoice")
		return fiber.ErrInternalServerError
	}

	return nil
}
