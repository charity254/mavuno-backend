CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL REFERENCES listings(id),
    sender_id UUID NOT NULL REFERENCES users(id),
    receiver_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_messages_listing_id ON messages(listing_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);