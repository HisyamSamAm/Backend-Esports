package handler

import (
	"embeck/model"
	"embeck/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllTeams godoc
// @Summary Get All Teams
// @Description Mendapatkan daftar semua tim dengan detail kapten
// @Tags Teams
// @Accept json
// @Produce json
// @Success 200 {array} model.TeamWithDetails
// @Failure 500 {object} map[string]interface{}
// @Router /api/admin/teams [get]
func GetAllTeams(c *fiber.Ctx) error {
	teams, err := repository.GetAllTeamsWithDetails(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data teams dari database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(teams)
}

// GetTeamByID godoc
// @Summary Get Team By ID
// @Description Mendapatkan detail tim berdasarkan ID dengan detail kapten
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} model.TeamWithDetails
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/admin/teams/{id} [get]
func GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := repository.GetTeamByIDWithDetails(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if team == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(team)
}

// CreateTeam godoc
// @Summary Create New Team
// @Description Membuat tim baru dengan kapten dan anggota
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body model.TeamRequest true "Team data"
// @Success 201 {object} model.TeamResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/admin/teams [post]
func CreateTeam(c *fiber.Ctx) error {
	var req model.TeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Validation
	if req.TeamName == "" || req.CaptainID == "" || len(req.Members) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "team_name, captain_id, and members are required",
		})
	}

	// Convert captain_id string to ObjectID
	captainObjID, err := primitive.ObjectIDFromHex(req.CaptainID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid captain_id format",
		})
	}

	// Convert members string array to ObjectID array
	var membersObjID []primitive.ObjectID
	for _, memberID := range req.Members {
		objID, err := primitive.ObjectIDFromHex(memberID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid member ID format: %s", memberID),
			})
		}
		membersObjID = append(membersObjID, objID)
	}

	// Validate captain is included in members
	captainInMembers := false
	for _, memberID := range membersObjID {
		if memberID == captainObjID {
			captainInMembers = true
			break
		}
	}
	if !captainInMembers {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Captain must be included in members list",
		})
	}

	// Create team model
	team := model.Team{
		TeamName:  req.TeamName,
		CaptainID: captainObjID,
		Members:   membersObjID,
		LogoURL:   req.LogoURL,
	}

	insertedID, err := repository.InsertTeam(c.Context(), team)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal menambahkan team: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Team created successfully",
		"team_id": insertedID,
	})
}

// UpdateTeam godoc
// @Summary Update Team
// @Description Memperbarui detail tim
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param request body model.TeamRequest true "Team data"
// @Success 200 {object} model.TeamResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/admin/teams/{id} [put]
func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.TeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Create update model (only update non-empty fields)
	update := model.Team{}

	if req.TeamName != "" {
		update.TeamName = req.TeamName
	}

	if req.CaptainID != "" {
		captainObjID, err := primitive.ObjectIDFromHex(req.CaptainID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid captain_id format",
			})
		}
		update.CaptainID = captainObjID
	}

	if len(req.Members) > 0 {
		var membersObjID []primitive.ObjectID
		for _, memberID := range req.Members {
			objID, err := primitive.ObjectIDFromHex(memberID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid member ID format: %s", memberID),
				})
			}
			membersObjID = append(membersObjID, objID)
		}
		update.Members = membersObjID

		// If both captain and members are provided, validate captain is in members
		if req.CaptainID != "" {
			captainInMembers := false
			for _, memberID := range membersObjID {
				if memberID == update.CaptainID {
					captainInMembers = true
					break
				}
			}
			if !captainInMembers {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Captain must be included in members list",
				})
			}
		}
	}

	if req.LogoURL != "" {
		update.LogoURL = req.LogoURL
	}

	_, err := repository.UpdateTeam(c.Context(), id, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Error updating team %s: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team updated successfully",
	})
}

// DeleteTeam godoc
// @Summary Delete Team
// @Description Menghapus tim dari database
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} model.TeamResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/admin/teams/{id} [delete]
func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeleteTeam(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Team dengan ID %s tidak ditemukan: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team deleted successfully",
	})
}
