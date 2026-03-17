CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    unit_type TEXT NOT NULL,
    description TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_farmer_id ON products(farmer_id);

CREATE INDEX idx_products_deleted ON products(deleted);