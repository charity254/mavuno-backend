-- Creates the produce_entries table.
-- This is the core feature of Mavuno — farmers record their daily
-- produce movements here: how much they had, sold, rejected etc.
CREATE TABLE produce_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL REFERENCES products(id),
    entry_date DATE NOT NULL,
    opening_stock INTEGER NOT NULL CHECK (opening_stock >= 0),
    added_stock INTEGER NOT NULL DEFAULT 0 CHECK (added_stock >= 0),
    sold_quantity INTEGER NOT NULL DEFAULT 0 CHECK (sold_quantity >= 0),
    rejected_quantity INTEGER NOT NULL DEFAULT 0 CHECK (rejected_quantity >= 0),
    price_per_unit INTEGER NOT NULL CHECK (price_per_unit >= 0),
    notes TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for faster queries when fetching a farmer's entries
CREATE INDEX idx_produce_entries_farmer_id ON produce_entries(farmer_id);

-- Index for faster queries when filtering by product
CREATE INDEX idx_produce_entries_product_id ON produce_entries(product_id);

-- Index for faster queries when filtering by date
CREATE INDEX idx_produce_entries_entry_date ON produce_entries(entry_date);