package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Alias   string             `bson:"alias" json:"alias"`
	LogoURL string             `bson:"logo_url" json:"logo_url"`
}

type Player struct {
	ID         string `bson:"id" json:"id"` // pakai "id", bukan "_id"
	Name       string `bson:"name" json:"name"`
	InGameName string `bson:"in_game_name" json:"in_game_name"`
	Role       string `bson:"role" json:"role"`
	TeamID     string `bson:"team_id" json:"team_id"` // juga pakai string
}


type PlayerInput struct {
	Name       string `json:"name"`
	InGameName string `json:"in_game_name"`
	Role       string `json:"role"`
	TeamID     string `json:"team_id"` // string dulu, nanti konversi
}
type Tournament struct {
	ObjectID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID         string             `bson:"id" json:"id"` // Custom string ID seperti "mpl-id-s14"
	Name       string             `bson:"name" json:"name"`
	Location   string             `bson:"location" json:"location"`
	StartDate  string             `bson:"start_date" json:"start_date"`
	EndDate    string             `bson:"end_date" json:"end_date"`
	Teams      []string           `bson:"teams" json:"teams"`
}


type Score struct {
	ID           string `bson:"id" json:"id"`                         // "match-1"
	TournamentID string `bson:"tournament_id" json:"tournament_id"`   // "mpl-id-s14"
	Team1ID      string `bson:"team1_id" json:"team1_id"`             // "team-fnoc"
	Team2ID      string `bson:"team2_id" json:"team2_id"`             // "team-tlid"
	Team1Score   int    `bson:"team1_score" json:"team1_score"`
	Team2Score   int    `bson:"team2_score" json:"team2_score"`
	Date         string `bson:"date" json:"date"`                     // karena di Mongo bentuknya string "2024-08-09"
}
