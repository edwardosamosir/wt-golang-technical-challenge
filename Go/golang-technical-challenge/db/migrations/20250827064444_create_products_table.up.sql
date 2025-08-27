CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id            UUID NOT NULL DEFAULT uuid_generate_v4(),
    invoice_no    VARCHAR(50) NOT NULL REFERENCES invoices(invoice_no) ON DELETE CASCADE,
    item_name     VARCHAR(255) NOT NULL CHECK (char_length(item_name) >= 5),
    quantity      INT NOT NULL CHECK (quantity >= 1),
    total_cost    DECIMAL(12,2) NOT NULL CHECK (total_cost >= 0),
    total_price   DECIMAL(12,2) NOT NULL CHECK (total_price >= 0),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_products_invoice_no 
    ON products(invoice_no);

COMMIT;