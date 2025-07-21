package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserTicket represents a ticket purchased by a user for a specific match.
type UserTicket struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	MatchID      primitive.ObjectID `bson:"match_id" json:"match_id"`
	PurchaseDate time.Time          `bson:"purchase_date" json:"purchase_date"`
	Status       string             `bson:"status" json:"status"` // e.g., "valid", "used"
}

// UserTicketRequest represents the request body for purchasing a ticket.
type UserTicketRequest struct {
	MatchID string `json:"match_id" validate:"required"`
}

// UserTicketResponse represents a single purchased ticket with populated match details.
type UserTicketResponse struct {
	ID           primitive.ObjectID `json:"_id"`
	UserID       primitive.ObjectID `json:"user_id"`
	MatchID      primitive.ObjectID `json:"match_id"`
	PurchaseDate time.Time          `json:"purchase_date"`
	Status       string             `json:"status"`
	MatchDetails *MatchBasicInfo    `json:"match_details,omitempty"`
}
