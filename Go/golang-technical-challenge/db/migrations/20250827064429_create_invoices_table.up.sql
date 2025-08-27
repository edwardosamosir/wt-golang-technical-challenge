BEGIN;

CREATE TYPE payment_enum AS ENUM ('CASH', 'CREDIT');

CREATE TABLE IF NOT EXISTS invoices (
    invoice_no       VARCHAR(50) PRIMARY KEY,
    date             DATE NOT NULL,
    customer_name    VARCHAR(255) NOT NULL CHECK (char_length(customer_name) >= 2),
    salesperson_name VARCHAR(255) NOT NULL CHECK (char_length(salesperson_name) >= 2),
    payment_type     payment_enum NOT NULL,
    notes            TEXT CHECK (notes IS NULL OR char_length(notes) >= 5),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;