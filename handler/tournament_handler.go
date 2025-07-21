package handler

import (
	"embeck/model"
	"embeck/repository"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateTournament creates a new tournament (admin only)
// @Summary Create new tournament
// @Description Create a new tournament (admin access required)
// @Tags Tournament Management (Admin)
// @Accept json
// @Produce json
// @Param tournament body model.TournamentRequest true "Tournament data"
// @Success 201 {object} model.TournamentResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/tournaments [post]
// @Security BearerAuth
func CreateTournament(c *fiber.Ctx) error {
	var req model.TournamentRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Description == "" || req.PrizePool == "" || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "missing_fields",
			Message: "Name, description, prize_pool, and status are required",
		})
	}

	// Validate status
	validStatuses := map[string]bool{
		"upcoming":  true,
		"ongoing":   true,
		"completed": true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_status",
			Message: "Status must be 'upcoming', 'ongoing', or 'completed'",
		})
	}

	// Validate date range
	if req.EndDate.Before(req.StartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_date_range",
			Message: "End date must be after start date",
		})
	}

	// Convert team IDs to ObjectIDs if provided
	var teamsParticipating []primitive.ObjectID
	if len(req.TeamsParticipating) > 0 {
		// Validate teams exist
		if err := repository.ValidateTeamsExist(req.TeamsParticipating); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "teams_not_found",
				Message: "One or more teams not found",
			})
		}

		for _, teamID := range req.TeamsParticipating {
			objID, err := primitive.ObjectIDFromHex(teamID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
					Error:   "invalid_team_id",
					Message: fmt.Sprintf("Invalid team ID: %s", teamID),
				})
			}
			teamsParticipating = append(teamsParticipating, objID)
		}
	}

	// Create tournament object
	tournament := model.Tournament{
		Name:               req.Name,
		Description:        req.Description,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		PrizePool:          req.PrizePool,
		RulesDocumentURL:   req.RulesDocumentURL,
		Status:             req.Status,
		TeamsParticipating: teamsParticipating,
		CreatedBy:          primitive.NewObjectID(), // TODO: Get from JWT
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save to database
	result, err := repository.CreateTournament(&tournament)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create tournament",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.TournamentResponse{
		Message:      "Tournament created successfully",
		TournamentID: result.InsertedID.(primitive.ObjectID).Hex(),
	})
}

// GetAllTournaments gets all tournaments (admin view)
// @Summary Get all tournaments
// @Description Get all tournaments with admin details
// @Tags Tournament Management (Admin)
// @Produce json
// @Success 200 {array} model.TournamentWithDetails
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/tournaments [get]
// @Security BearerAuth
func GetAllTournaments(c *fiber.Ctx) error {
	tournaments, err := repository.GetAllTournaments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve tournaments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournaments)
}

// GetTournamentByID gets tournament by ID (admin view)
// @Summary Get tournament by ID
// @Description Get tournament details by ID (admin access)
// @Tags Tournament Management (Admin)
// @Produce json
// @Param id path string true "Tournament ID"
// @Success 200 {object} model.Tournament
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/tournaments/{id} [get]
// @Security BearerAuth
func GetTournamentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// MEMASTIKAN FUNGSI INI MENGEMBALIKAN DETAIL LENGKAP
	tournament, err := repository.GetTournamentWithDetailsByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve tournament",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}

// UpdateTournament updates tournament data
// @Summary Update tournament
// @Description Update tournament details including participating teams
// @Tags Tournament Management (Admin)
// @Accept json
// @Produce json
// @Param id path string true "Tournament ID"
// @Param tournament body model.TournamentRequest true "Tournament data"
// @Success 200 {object} model.TournamentResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/tournaments/{id} [put]
// @Security BearerAuth
func UpdateTournament(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.TournamentRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := map[string]bool{
			"upcoming":  true,
			"ongoing":   true,
			"completed": true,
		}
		if !validStatuses[req.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "invalid_status",
				Message: "Status must be 'upcoming', 'ongoing', or 'completed'",
			})
		}
	}

	// Validate date range if both dates provided
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() && req.EndDate.Before(req.StartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_date_range",
			Message: "End date must be after start date",
		})
	}

	// Build update document
	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if !req.StartDate.IsZero() {
		update["start_date"] = req.StartDate
	}
	if !req.EndDate.IsZero() {
		update["end_date"] = req.EndDate
	}
	if req.PrizePool != "" {
		update["prize_pool"] = req.PrizePool
	}
	if req.RulesDocumentURL != "" {
		update["rules_document_url"] = req.RulesDocumentURL
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	// Handle teams participating
	if len(req.TeamsParticipating) > 0 {
		// Validate teams exist
		if err := repository.ValidateTeamsExist(req.TeamsParticipating); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "teams_not_found",
				Message: "One or more teams not found",
			})
		}

		var teamsParticipating []primitive.ObjectID
		for _, teamID := range req.TeamsParticipating {
			objID, err := primitive.ObjectIDFromHex(teamID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
					Error:   "invalid_team_id",
					Message: fmt.Sprintf("Invalid team ID: %s", teamID),
				})
			}
			teamsParticipating = append(teamsParticipating, objID)
		}
		update["teams_participating"] = teamsParticipating
	}

	// Update tournament
	err := repository.UpdateTournament(id, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update tournament",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.TournamentResponse{
		Message: "Tournament updated successfully",
	})
}

// DeleteTournament deletes a tournament
// @Summary Delete tournament
// @Description Delete tournament by ID
// @Tags Tournament Management (Admin)
// @Produce json
// @Param id path string true "Tournament ID"
// @Success 200 {object} model.TournamentResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/tournaments/{id} [delete]
// @Security BearerAuth
func DeleteTournament(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeleteTournament(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.TournamentResponse{
		Message: "Tournament deleted successfully",
	})
}

// GetAllTournamentsPublic gets all tournaments for public access
// @Summary Get all tournaments (public)
// @Description Get all tournaments without admin details
// @Tags Tournament Data (Public)
// @Produce json
// @Success 200 {array} model.TournamentPublic
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tournaments [get]
func GetAllTournamentsPublic(c *fiber.Ctx) error {
	tournaments, err := repository.GetAllTournamentsPublic()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve tournaments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournaments)
}

// GetTournamentWithDetailsByID gets tournament with populated details (public)
// @Summary Get tournament details (public)
// @Description Get tournament with populated teams and matches
// @Tags Tournament Data (Public)
// @Produce json
// @Param id path string true "Tournament ID"
// @Success 200 {object} model.TournamentWithDetails
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tournaments/{id} [get]
func GetTournamentWithDetailsByID(c *fiber.Ctx) error {
	id := c.Params("id")

	tournament, err := repository.GetTournamentWithDetailsByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve tournament",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
