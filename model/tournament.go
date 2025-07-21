package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tournament represents a tournament entity
type Tournament struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name               string               `bson:"name" json:"name"`
	Description        string               `bson:"description" json:"description"`
	StartDate          time.Time            `bson:"start_date" json:"start_date"`
	EndDate            time.Time            `bson:"end_date" json:"end_date"`
	PrizePool          string               `bson:"prize_pool" json:"prize_pool"`
	RulesDocumentURL   string               `bson:"rules_document_url,omitempty" json:"rules_document_url,omitempty"`
	Status             string               `bson:"status" json:"status"`
	TeamsParticipating []primitive.ObjectID `bson:"teams_participating" json:"teams_participating"`
	CreatedBy          primitive.ObjectID   `bson:"created_by" json:"created_by"`
	CreatedAt          time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time            `bson:"updated_at" json:"updated_at"`
}

// TournamentRequest represents request body for creating/updating tournament
type TournamentRequest struct {
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description" validate:"required"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	PrizePool          string    `json:"prize_pool" validate:"required"`
	RulesDocumentURL   string    `json:"rules_document_url,omitempty"`
	Status             string    `json:"status" validate:"required,oneof=upcoming ongoing completed"`
	TeamsParticipating []string  `json:"teams_participating,omitempty"`
}

// TournamentResponse represents response for tournament operations
type TournamentResponse struct {
	Message      string `json:"message"`
	TournamentID string `json:"tournament_id,omitempty"`
}

// TournamentPublic represents tournament data for public access (without admin fields)
type TournamentPublic struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	StartDate        time.Time          `bson:"start_date" json:"start_date"`
	EndDate          time.Time          `bson:"end_date" json:"end_date"`
	PrizePool        string             `bson:"prize_pool" json:"prize_pool"`
	RulesDocumentURL string             `bson:"rules_document_url,omitempty" json:"rules_document_url,omitempty"`
	Status           string             `bson:"status" json:"status"`
}

// TournamentWithDetails represents tournament with populated teams and matches
type TournamentWithDetails struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	StartDate          time.Time          `bson:"start_date" json:"start_date"`
	EndDate            time.Time          `bson:"end_date" json:"end_date"`
	PrizePool          string             `bson:"prize_pool" json:"prize_pool"`
	RulesDocumentURL   string             `bson:"rules_document_url,omitempty" json:"rules_document_url,omitempty"`
	Status             string             `bson:"status" json:"status"`
	TeamsParticipating []TeamBasicInfo    `bson:"teams_participating,omitempty" json:"teams_participating"`
	Matches            []MatchBasicInfo   `bson:"matches,omitempty" json:"matches"`
}

// TeamBasicInfo represents minimal team info for tournament details
type TeamBasicInfo struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	TeamName string             `json:"team_name" bson:"team_name"`
	LogoURL  string             `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
}

// MatchBasicInfo represents minimal match info for tournament details
type MatchBasicInfo struct {
	ID               primitive.ObjectID  `bson:"_id" json:"_id"`
	MatchDate        time.Time           `bson:"match_date" json:"match_date"`
	MatchTime        string              `bson:"match_time" json:"match_time"`
	Location         string              `bson:"location" json:"location"`
	Round            string              `bson:"round" json:"round"`
	TeamA            TeamBasicInfo       `bson:"team_a" json:"team_a"`
	TeamB            TeamBasicInfo       `bson:"team_b" json:"team_b"`
	ResultTeamAScore *int                `bson:"result_team_a_score,omitempty" json:"result_team_a_score"`
	ResultTeamBScore *int                `bson:"result_team_b_score,omitempty" json:"result_team_b_score"`
	WinnerTeamID     *primitive.ObjectID `bson:"winner_team_id,omitempty" json:"winner_team_id"`
	Status           string              `bson:"status" json:"status"`
}
