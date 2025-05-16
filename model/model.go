package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Alias   string             `bson:"alias" json:"alias"`
	LogoURL string             `bson:"logo_url" json:"logo_url"`
}

type Player struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	InGameName string             `bson:"in_game_name" json:"in_game_name"`
	Role       string             `bson:"role" json:"role"`
	TeamID     primitive.ObjectID `bson:"team_id" json:"team_id"`
}

type Tournament struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Season      string             `bson:"season" json:"season"`
	StartDate   string             `bson:"start_date" json:"start_date"`
	EndDate     string             `bson:"end_date" json:"end_date"`
	Description string             `bson:"description" json:"description"`
}

type Score struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
    Team1ID      primitive.ObjectID `bson:"team1_id" json:"team1_id"`
    Team2ID      primitive.ObjectID `bson:"team2_id" json:"team2_id"`
    Team1Score   int                `bson:"team1_score" json:"team1_score"`
    Team2Score   int                `bson:"team2_score" json:"team2_score"`
    Date         time.Time          `bson:"date" json:"date"`
}