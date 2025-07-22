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
	TeamName  string   `json:"team_name" bson:"team_name" example:"RRQ Hoshi"`
	CaptainID string   `json:"captain_id" bson:"captain_id" example:"687f9d7c8efa8f58af86646a"`
	Members   []string `json:"members" bson:"members" example:"687f9d7c8efa8f58af86646a,687f9d7c8efa8f58af86646b"`
	LogoURL   string   `json:"logo_url,omitempty" bson:"logo_url,omitempty" example:"https://example.com/rrq_logo.png"`
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
	CaptainDetails *PlayerDetails       `json:"captain_details,omitempty" bson:"captain_details,omitempty"`
	MembersDetails []PlayerDetails      `json:"members_details,omitempty" bson:"members_details,omitempty"`
}

// PlayerDetails represents minimal player info for team details
type PlayerDetails struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	MLNickname string             `json:"ml_nickname" bson:"ml_nickname"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	MLID       string             `json:"ml_id,omitempty" bson:"ml_id,omitempty"`
	Status     string             `json:"status,omitempty" bson:"status,omitempty"`
}
