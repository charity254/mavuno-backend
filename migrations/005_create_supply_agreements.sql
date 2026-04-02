CREATE TABLE supply_agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL REFERENCES products(id),
    supply_location_id UUID NOT NULL REFERENCES supply_locations(id),
    quantity_per_delivery INTEGER NOT NULL CHECK (quantity_per_delivery > 0),
    price_per_unit INTEGER NOT NULL CHECK (price_per_unit >= 0),
    delivery_days TEXT[] NOT NULL,  -- array of weekday names e.g. ["Monday", "Wednesday"]
    delivery_notes TEXT,
    active BOOLEAN DEFAULT TRUE,    -- can be toggled on/off without deleting
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_supply_agreements_farmer_id ON supply_agreements(farmer_id);
CREATE INDEX idx_supply_agreements_product_id ON supply_agreements(product_id);
CREATE INDEX idx_supply_agreements_location_id ON supply_agreements(supply_location_id);