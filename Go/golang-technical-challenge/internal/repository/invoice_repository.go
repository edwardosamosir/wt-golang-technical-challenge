package repository

import (
	"golang-technical-challenge/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	Repository[entity.Invoice]
	Log *logrus.Logger
}

func NewInvoiceRepository(log *logrus.Logger) *InvoiceRepository {
	return &InvoiceRepository{
		Log: log,
	}
}

func (r *InvoiceRepository) FindByInvoiceNo(db *gorm.DB, invoice *entity.Invoice, invoiceNo string) error {
	return db.Preload("Products").
		Where("invoice_no = ?", invoiceNo).
		Take(invoice).Error
}

func (r *InvoiceRepository) FindInvoicesByNumbers(db *gorm.DB, invoiceNos []string) ([]entity.Invoice, error) {
	if len(invoiceNos) == 0 {
		return []entity.Invoice{}, nil
	}

	var invoices []entity.Invoice
	if err := db.Preload("Products").
		Where("invoice_no IN ?", invoiceNos).
		Order("date DESC, created_at DESC").
		Find(&invoices).Error; err != nil {
		r.Log.WithError(err).
			WithField("invoice_nos", invoiceNos).
			Error("Failed to find invoices by invoice numbers")
		return nil, err
	}

	return invoices, nil
}

func (r *InvoiceRepository) FindInvoicesByDate(db *gorm.DB, date string, limit, offset int) ([]entity.Invoice, int64, error) {
	var invoices []entity.Invoice
	var total int64

	query := db.Where("date = ?", date)

	if err := query.Model(&entity.Invoice{}).Count(&total).Error; err != nil {
		r.Log.WithError(err).
			WithField("date", date).
			Error("Failed to count total invoices by date")
		return nil, 0, err
	}

	if err := query.Preload("Products").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error; err != nil {
		r.Log.WithError(err).
			WithFields(logrus.Fields{
				"date":   date,
				"limit":  limit,
				"offset": offset,
			}).
			Error("Failed to find invoices by date")
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *InvoiceRepository) GetSummaryByDate(db *gorm.DB, date string) (totalProfit, totalCash string, err error) {
	type result struct {
		TotalProfit string
		TotalCash   string
	}

	var res result
	query := `
		SELECT 
			COALESCE(SUM((p.total_price - p.total_cost) * p.quantity), 0)::text AS total_profit,
			COALESCE(SUM(
				CASE 
					WHEN i.payment_type = 'CASH' 
					THEN (p.total_price * p.quantity) 
					ELSE 0 
				END
			), 0)::text AS total_cash

		FROM products p
		JOIN invoices i ON i.invoice_no = p.invoice_no
		WHERE i.date = ?
	`

	if err := db.Raw(query, date).Scan(&res).Error; err != nil {
		r.Log.WithError(err).WithField("date", date).Error("Failed to calculate invoice summary")
		return "0", "0", err
	}

	return res.TotalProfit, res.TotalCash, nil
}
