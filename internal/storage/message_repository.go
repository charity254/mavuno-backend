package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(m *models.Message) error {
	query := `
		INSERT INTO messages (id, listing_id, sender_id, receiver_id, content, created_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query,
		m.ID, m.ListingID, m.SenderID,
		m.ReceiverID, m.Content, m.CreatedAt, m.Deleted,
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetMessagesByListing(listingID uuid.UUID) ([]*models.Message, error) {
	query := `
		SELECT id, listing_id, sender_id, receiver_id, content, created_at, deleted
		FROM messages
		WHERE listing_id = $1 AND deleted = false
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		err := rows.Scan(
			&m.ID, &m.ListingID, &m.SenderID,
			&m.ReceiverID, &m.Content, &m.CreatedAt, &m.Deleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepository) IsParticipant(listingID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM messages
		WHERE listing_id = $1
		AND (sender_id = $2 OR receiver_id = $2)
		AND deleted = false
	`
	var count int
	err := r.db.QueryRow(query, listingID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check participant: %w", err)
	}
	return count > 0, nil
}