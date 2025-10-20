CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idn TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_customers_idn ON customers(idn);

COMMENT ON TABLE customers IS 'Stores customer information with unique IDN (ИИН/БИН)';
COMMENT ON COLUMN customers.idn IS 'ИИН/БИН - 12 digit unique identifier';
