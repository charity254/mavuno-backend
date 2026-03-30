
-- A supply location is a place a farmer regularly delivers produce to
CREATE TABLE supply_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,                
    contact_person TEXT NOT NULL,      
    phone_number TEXT NOT NULL,        
    location_address TEXT NOT NULL,    
    notes TEXT,                        
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for faster queries when fetching a farmer's supply locations
CREATE INDEX idx_supply_locations_farmer_id ON supply_locations(farmer_id);