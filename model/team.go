package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Team represents a Mobile Legends team
type Team struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	TeamName  string               `bson:"team_name" json:"team_name"`
	CaptainID primitive.ObjectID   `bson:"captain_id" json:"captain_id"`
	Members   []primitive.ObjectID `bson:"members" json:"members"`
	LogoURL   string               `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

// TeamRequest represents request body for creating/updating team
type TeamRequest struct {
	TeamName  string   `json:"team_name" bson:"team_name"`
	CaptainID string   `json:"captain_id" bson:"captain_id"`
	Members   []string `json:"members" bson:"members"`
	LogoURL   string   `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
}

// TeamResponse represents response for team operations
type TeamResponse struct {
	Message string `json:"message"`
	TeamID  string `json:"team_id,omitempty"`
}

// TeamWithDetails represents team with populated captain details
type TeamWithDetails struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	TeamName       string               `bson:"team_name" json:"team_name"`
	CaptainID      primitive.ObjectID   `bson:"captain_id" json:"captain_id"`
	Members        []primitive.ObjectID `bson:"members" json:"members"`
	LogoURL        string               `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
	CaptainDetails *PlayerDetails       `json:"captain_details,omitempty"`
}

// PlayerDetails represents minimal player info for team details
type PlayerDetails struct {
	ID         primitive.ObjectID `json:"_id"`
	MLNickname string             `json:"ml_nickname"`
}
