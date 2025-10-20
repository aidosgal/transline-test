CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route TEXT NOT NULL,
    price NUMERIC NOT NULL CHECK (price >= 0),
    status TEXT NOT NULL DEFAULT 'CREATED',
    customer_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shipments_customer_id ON shipments(customer_id);
CREATE INDEX IF NOT EXISTS idx_shipments_status ON shipments(status);
CREATE INDEX IF NOT EXISTS idx_shipments_created_at ON shipments(created_at DESC);

COMMENT ON TABLE shipments IS 'Stores shipment/delivery orders';
COMMENT ON COLUMN shipments.customer_id IS 'Reference to customer UUID from customer-service';
COMMENT ON COLUMN shipments.status IS 'Shipment status: CREATED, IN_PROGRESS, DELIVERED, etc.';
