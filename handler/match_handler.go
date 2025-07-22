package handler

import (
	"embeck/model"
	"embeck/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateMatch godoc
// @Summary Create New Match
// @Description Membuat pertandingan baru untuk turnamen
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.MatchRequest true "Match data"
// @Success 201 {object} model.MatchResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /api/admin/matches [post]
func CreateMatch(c *fiber.Ctx) error {
	var req model.MatchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request data",
		})
	}

	// Validation
	if req.TournamentID == "" || req.TeamAID == "" || req.TeamBID == "" || req.MatchTime == "" || req.Round == "" || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "missing_fields",
			Message: "tournament_id, team_a_id, team_b_id, match_date, match_time, round, and status are required",
		})
	}

	// Validate status
	validStatuses := map[string]bool{
		"scheduled": true,
		"ongoing":   true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_status",
			Message: "Status must be 'scheduled', 'ongoing', 'completed', or 'cancelled'",
		})
	}

	// Convert string IDs to ObjectIDs
	tournamentObjID, err := primitive.ObjectIDFromHex(req.TournamentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid tournament_id format",
		})
	}

	teamAObjID, err := primitive.ObjectIDFromHex(req.TeamAID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid team_a_id format",
		})
	}

	teamBObjID, err := primitive.ObjectIDFromHex(req.TeamBID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid team_b_id format",
		})
	}

	// Validate teams are different
	if req.TeamAID == req.TeamBID {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Team A and Team B must be different",
		})
	}

	// Create match model
	match := model.Match{
		TournamentID:     tournamentObjID,
		TeamAID:          teamAObjID,
		TeamBID:          teamBObjID,
		MatchDate:        req.MatchDate,
		MatchTime:        req.MatchTime,
		Location:         req.Location,
		Round:            req.Round,
		ResultTeamAScore: req.ResultTeamAScore,
		ResultTeamBScore: req.ResultTeamBScore,
		Status:           req.Status,
	}

	// Handle winner team ID if provided
	if req.WinnerTeamID != "" {
		winnerObjID, err := primitive.ObjectIDFromHex(req.WinnerTeamID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "invalid_id",
				Message: "Invalid winner_team_id format",
			})
		}
		match.WinnerTeamID = &winnerObjID

		// Validate winner is either team A or team B
		if req.WinnerTeamID != req.TeamAID && req.WinnerTeamID != req.TeamBID {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "validation_error",
				Message: "Winner team must be either team A or team B",
			})
		}
	}

	insertedID, err := repository.CreateMatch(c.Context(), match)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
			Error:   "db_conflict",
			Message: fmt.Sprintf("Gagal menambahkan match: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.MatchResponse{
		Message: "Match created successfully",
		MatchID: insertedID.(primitive.ObjectID).Hex(),
	})
}

// GetAllMatches godoc
// @Summary Get All Matches
// @Description Mendapatkan daftar semua pertandingan, bisa difilter berdasarkan tournament_id
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tournament_id query string false "Filter pertandingan berdasarkan ID turnamen"
// @Success 200 {array} model.MatchWithDetails
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/matches [get]
func GetAllMatches(c *fiber.Ctx) error {
	// Get optional tournament_id filter from query params
	tournamentID := c.Query("tournament_id")

	matches, err := repository.GetAllMatches(c.Context(), tournamentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "db_error",
			Message: "Gagal mengambil data matches dari database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(matches)
}

// GetMatchByID godoc
// @Summary Get Match By ID
// @Description Mendapatkan detail pertandingan berdasarkan ID
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Success 200 {object} model.Match
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/matches/{id} [get]
func GetMatchByID(c *fiber.Ctx) error {
	id := c.Params("id")

	match, err := repository.GetMatchWithDetailsByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_id",
			Message: err.Error(),
		})
	}

	if match == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: "Match not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(match)
}

// UpdateMatch godoc
// @Summary Update Match
// @Description Memperbarui detail pertandingan termasuk input skor
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Param request body model.MatchRequest true "Match data"
// @Success 200 {object} model.MatchResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/matches/{id} [put]
func UpdateMatch(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.MatchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request data",
		})
	}

	update := bson.M{}

	if req.TournamentID != "" {
		tournamentObjID, err := primitive.ObjectIDFromHex(req.TournamentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_id", Message: "Invalid tournament_id format"})
		}
		update["tournament_id"] = tournamentObjID
	}

	if req.TeamAID != "" {
		teamAObjID, err := primitive.ObjectIDFromHex(req.TeamAID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_id", Message: "Invalid team_a_id format"})
		}
		update["team_a_id"] = teamAObjID
	}

	if req.TeamBID != "" {
		teamBObjID, err := primitive.ObjectIDFromHex(req.TeamBID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_id", Message: "Invalid team_b_id format"})
		}
		update["team_b_id"] = teamBObjID
	}

	// Validate teams are different if both are provided
	if req.TeamAID != "" && req.TeamBID != "" && req.TeamAID == req.TeamBID {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "validation_error", Message: "Team A and Team B must be different"})
	}

	if !req.MatchDate.IsZero() {
		update["match_date"] = req.MatchDate
	}

	if req.MatchTime != "" {
		update["match_time"] = req.MatchTime
	}

	if req.Location != "" {
		update["location"] = req.Location
	}

	if req.Round != "" {
		update["round"] = req.Round
	}

	if req.ResultTeamAScore != nil {
		update["result_team_a_score"] = req.ResultTeamAScore
	}

	if req.ResultTeamBScore != nil {
		update["result_team_b_score"] = req.ResultTeamBScore
	}

	if req.WinnerTeamID != "" {
		winnerObjID, err := primitive.ObjectIDFromHex(req.WinnerTeamID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_id", Message: "Invalid winner_team_id format"})
		}
		update["winner_team_id"] = &winnerObjID
	}

	if req.Status != "" {
		validStatuses := map[string]bool{"scheduled": true, "ongoing": true, "completed": true, "cancelled": true}
		if !validStatuses[req.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_status", Message: "Status must be 'scheduled', 'ongoing', 'completed', or 'cancelled'"})
		}
		update["status"] = req.Status
	}

	if len(update) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "bad_request", Message: "No fields to update"})
	}

	_, err := repository.UpdateMatch(c.Context(), id, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "update_failed",
			Message: fmt.Sprintf("Error updating match %s: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.MatchResponse{
		Message: "Match updated successfully",
	})
}

// DeleteMatch godoc
// @Summary Delete Match
// @Description Menghapus pertandingan dari database
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Success 200 {object} model.MatchResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/matches/{id} [delete]
func DeleteMatch(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeleteMatch(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "not_found",
			Message: fmt.Sprintf("Match dengan ID %s tidak ditemukan: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.MatchResponse{
		Message: "Match deleted successfully",
	})
}
