package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID         string          `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	InvoiceNo  string          `gorm:"column:invoice_no;type:varchar(50);not null;index"`
	ItemName   string          `gorm:"column:item_name;type:varchar(255);not null;check:char_length(item_name) >= 5"`
	Quantity   int             `gorm:"column:quantity;not null;check:quantity >= 1"`
	TotalCost  decimal.Decimal `gorm:"column:total_cost;type:decimal(12,2);not null;check:total_cost >= 0"`
	TotalPrice decimal.Decimal `gorm:"column:total_price;type:decimal(12,2);not null;check:total_price >= 0"`
	CreatedAt  time.Time       `gorm:"column:created_at;type:timestamptz;default:now();not null"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;type:timestamptz;default:now()"`

	Invoice Invoice `gorm:"foreignKey:InvoiceNo;references:InvoiceNo;constraint:OnDelete:CASCADE"`
}

func (Product) TableName() string {
	return "products"
}
