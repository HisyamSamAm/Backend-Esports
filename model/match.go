package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Match represents a match entity
type Match struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	TournamentID     primitive.ObjectID  `bson:"tournament_id" json:"tournament_id"`
	TeamAID          primitive.ObjectID  `bson:"team_a_id" json:"team_a_id"`
	TeamBID          primitive.ObjectID  `bson:"team_b_id" json:"team_b_id"`
	MatchDate        time.Time           `bson:"match_date" json:"match_date"`
	MatchTime        string              `bson:"match_time" json:"match_time"`
	Location         string              `bson:"location,omitempty" json:"location,omitempty"`
	Round            string              `bson:"round" json:"round"`
	ResultTeamAScore *int                `bson:"result_team_a_score,omitempty" json:"result_team_a_score"`
	ResultTeamBScore *int                `bson:"result_team_b_score,omitempty" json:"result_team_b_score"`
	WinnerTeamID     *primitive.ObjectID `bson:"winner_team_id,omitempty" json:"winner_team_id,omitempty"`
	Status           string              `bson:"status" json:"status"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updated_at"`
}

// MatchRequest represents request body for creating/updating match
type MatchRequest struct {
	TournamentID     string    `json:"tournament_id" validate:"required"`
	TeamAID          string    `json:"team_a_id" validate:"required"`
	TeamBID          string    `json:"team_b_id" validate:"required"`
	MatchDate        time.Time `json:"match_date" validate:"required"`
	MatchTime        string    `json:"match_time" validate:"required"`
	Location         string    `json:"location,omitempty"`
	Round            string    `json:"round" validate:"required"`
	ResultTeamAScore *int      `json:"result_team_a_score,omitempty"`
	ResultTeamBScore *int      `json:"result_team_b_score,omitempty"`
	WinnerTeamID     string    `json:"winner_team_id,omitempty"`
	Status           string    `json:"status" validate:"required,oneof=scheduled ongoing completed cancelled"`
}

// MatchResponse represents response for match operations
type MatchResponse struct {
	Message string `json:"message"`
	MatchID string `json:"match_id,omitempty"`
}

// MatchWithDetails represents match with populated team details
type MatchWithDetails struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	TournamentID     primitive.ObjectID  `bson:"tournament_id" json:"tournament_id"`
	TeamAID          primitive.ObjectID  `bson:"team_a_id" json:"team_a_id"`
	TeamBID          primitive.ObjectID  `bson:"team_b_id" json:"team_b_id"`
	MatchDate        time.Time           `bson:"match_date" json:"match_date"`
	MatchTime        string              `bson:"match_time" json:"match_time"`
	Location         string              `bson:"location,omitempty" json:"location,omitempty"`
	Round            string              `bson:"round" json:"round"`
	ResultTeamAScore *int                `bson:"result_team_a_score,omitempty" json:"result_team_a_score"`
	ResultTeamBScore *int                `bson:"result_team_b_score,omitempty" json:"result_team_b_score"`
	WinnerTeamID     *primitive.ObjectID `bson:"winner_team_id,omitempty" json:"winner_team_id,omitempty"`
	Status           string              `bson:"status" json:"status"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updated_at"`
	TeamA            *TeamBasicInfo      `json:"team_a,omitempty" bson:"team_a,omitempty"`
	TeamB            *TeamBasicInfo      `json:"team_b,omitempty" bson:"team_b,omitempty"`
}
