CREATE TABLE listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id UUID NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    unit_type TEXT NOT NULL,
    location TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_listings_farmer_id ON listings(farmer_id);
CREATE INDEX idx_listings_deleted ON listings(deleted);