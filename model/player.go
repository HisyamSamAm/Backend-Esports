package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Player represents a Mobile Legends player
type Player struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	MLNickname string             `bson:"ml_nickname" json:"ml_nickname"`
	MLID       string             `bson:"ml_id" json:"ml_id"`
	Status     string             `bson:"status" json:"status"`
	AvatarURL  string             `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// PlayerRequest represents request body for creating/updating player
type PlayerRequest struct {
	Name       string `json:"name" bson:"name"`
	MLNickname string `json:"ml_nickname" bson:"ml_nickname"`
	MLID       string `json:"ml_id" bson:"ml_id"`
	Status     string `json:"status" bson:"status"`
	AvatarURL  string `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
}

// PlayerResponse represents response for player operations
type PlayerResponse struct {
	Message  string `json:"message"`
	PlayerID string `json:"player_id,omitempty"`
}
