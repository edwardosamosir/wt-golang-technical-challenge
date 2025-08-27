package entity

import "time"

type Invoice struct {
	InvoiceNo       string    `gorm:"column:invoice_no;type:varchar(50);primaryKey"`
	Date            time.Time `gorm:"column:date;type:date;not null"`
	CustomerName    string    `gorm:"column:customer_name;type:varchar(255);not null;check:char_length(customer_name) >= 2"`
	SalespersonName string    `gorm:"column:salesperson_name;type:varchar(255);not null;check:char_length(salesperson_name) >= 2"`
	PaymentType     string    `gorm:"column:payment_type;type:payment_enum;not null"`
	Notes           *string   `gorm:"column:notes;check:notes IS NULL OR char_length(notes) >= 5"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamptz;default:now();not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamptz;default:now();not null"`

	Products []Product `gorm:"foreignKey:InvoiceNo;references:InvoiceNo;constraint:OnDelete:CASCADE"`
}

func (Invoice) TableName() string {
	return "invoices"
}
