package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type MessageService struct {
	repo        *storage.MessageRepository
	listingRepo *storage.ListingRepository
}

func NewMessageService(repo *storage.MessageRepository, listingRepo *storage.ListingRepository) *MessageService {
	return &MessageService{
		repo:        repo,
		listingRepo: listingRepo,
	}
}

func (s *MessageService) SendMessage(listingID, senderID uuid.UUID, content string) (*models.Message, error) {
	if content == "" {
		return nil, fmt.Errorf("message content is required")
	}
	if len(content) > 1000 {
		return nil, fmt.Errorf("message cannot exceed 1000 characters")
	}

	// Fetch the listing to get the farmer ID as the receiver
	listing, err := s.listingRepo.GetListingByID(listingID)
	if err != nil {
		return nil, fmt.Errorf("listing not found")
	}

	// The receiver is always the farmer who owns the listing
	// unless the sender IS the farmer in which case we need
	// the other participant as the receiver
	receiverID := listing.FarmerID
	if senderID == listing.FarmerID {
		return nil, fmt.Errorf("farmers cannot message themselves on their own listing")
	}

	m := &models.Message{
		ID:         uuid.New(),
		ListingID:  listingID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
		Deleted:    false,
	}

	if err := s.repo.CreateMessage(m); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return m, nil
}

func (s *MessageService) GetMessages(listingID, userID uuid.UUID) ([]*models.Message, error) {
	// Fetch the listing to check if user is the farmer
	listing, err := s.listingRepo.GetListingByID(listingID)
	if err != nil {
		return nil, fmt.Errorf("listing not found")
	}

	// Allow access if user is the farmer who owns the listing
	if listing.FarmerID != userID {
		// Otherwise check if user is a participant in the conversation
		isParticipant, err := s.repo.IsParticipant(listingID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to verify access: %w", err)
		}
		if !isParticipant {
			return nil, fmt.Errorf("you do not have access to these messages")
		}
	}

	messages, err := s.repo.GetMessagesByListing(listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	if messages == nil {
		messages = []*models.Message{}
	}

	return messages, nil
}