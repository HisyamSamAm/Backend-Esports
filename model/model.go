package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"507f1f77bcf86cd799439011"`
	Name    string             `bson:"name" json:"name" example:"Fnatic ONIC"`
	Alias   string             `bson:"alias" json:"alias" example:"FNOC"`
	LogoURL string             `bson:"logo_url" json:"logo_url" example:"https://example.com/logo.png"`
}

type Player struct {
	ID         string `bson:"id" json:"id" example:"player-001"`
	Name       string `bson:"name" json:"name" example:"John Doe"`
	InGameName string `bson:"in_game_name" json:"in_game_name" example:"JohnMLBB"`
	Role       string `bson:"role" json:"role" example:"Tank"`
	TeamID     string `bson:"team_id" json:"team_id" example:"team-fnoc"`
}

type PlayerInput struct {
	Name       string `json:"name" example:"John Doe"`
	InGameName string `json:"in_game_name" example:"JohnMLBB"`
	Role       string `json:"role" example:"Tank"`
	TeamID     string `json:"team_id" example:"team-fnoc"`
}

type Tournament struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID        string             `bson:"id" json:"id"` // Custom string ID seperti "mpl-id-s14"
	Name      string             `bson:"name" json:"name"`
	Location  string             `bson:"location" json:"location"`
	StartDate string             `bson:"start_date" json:"start_date"`
	EndDate   string             `bson:"end_date" json:"end_date"`
	Teams     []string           `bson:"teams" json:"teams"`
}

type Score struct {
	ID           string `bson:"id" json:"id"`                       // "match-1"
	TournamentID string `bson:"tournament_id" json:"tournament_id"` // "mpl-id-s14"
	Team1ID      string `bson:"team1_id" json:"team1_id"`           // "team-fnoc"
	Team2ID      string `bson:"team2_id" json:"team2_id"`           // "team-tlid"
	Team1Score   int    `bson:"team1_score" json:"team1_score"`
	Team2Score   int    `bson:"team2_score" json:"team2_score"`
	Date         string `bson:"date" json:"date"` // karena di Mongo bentuknya string "2024-08-09"
}

// Response models for Swagger documentation
type APIResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Success message"`
	Data    interface{} `json:"data,omitempty"`
}

type TeamResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Success ambil 1 team!"`
	Data    Team   `json:"data"`
}

type TeamsResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"success ngambil data bre!"`
	Data    []Team `json:"data"`
}

type PlayerResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Berhasil ambil player!"`
	Data    Player `json:"data"`
}

type PlayersResponse struct {
	Status  int      `json:"status" example:"200"`
	Message string   `json:"message" example:"Berhasil ambil data players!"`
	Data    []Player `json:"data"`
}

type ErrorResponse struct {
	Status  int         `json:"status" example:"400"`
	Message string      `json:"message" example:"Error message"`
	Data    interface{} `json:"data"`
}
