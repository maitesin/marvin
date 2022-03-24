package domain

import "github.com/google/uuid"

type Channel struct {
	ID         uuid.UUID
	DeliveryID string
	ChatID     int64
}

func NewChannel(id uuid.UUID, deliveryID string, chatID int64) Channel {
	return Channel{
		ID:         id,
		DeliveryID: deliveryID,
		ChatID:     chatID,
	}
}
