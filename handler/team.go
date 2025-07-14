package handler

import (
	"EMBECK/model"
	repo "EMBECK/repository"

	"github.com/gofiber/fiber/v2"
)

// GetAllTeams godoc
// @Summary Get all teams
// @Description Get list of all teams
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} model.TeamsResponse "success"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /team [get]
func GetAllTeams(c *fiber.Ctx) error {
	team, err := repo.GetAllTeams(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error nih servernya bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success ngambil data bre!",
		"data":    team,
	})
}

// GetTeamByID godoc
// @Summary Get team by ID
// @Description Get a single team by its ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} model.TeamResponse "success"
// @Failure 404 {object} model.ErrorResponse "team not found"
// @Router /team/{id} [get]
func GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := repo.GetTeamByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Team gak ketemu bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success ambil 1 team!",
		"data":    team,
	})
}

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team with the provided data
// @Tags teams
// @Accept json
// @Produce json
// @Param team body model.Team true "Team data"
// @Success 201 {object} model.APIResponse "team created"
// @Failure 400 {object} model.ErrorResponse "invalid data"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /team [post]
func CreateTeam(c *fiber.Ctx) error {
	var team model.Team

	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	id, err := repo.CreateTeam(c.Context(), team)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Team berhasil ditambahkan bre!",
		"data":    id,
	})
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update an existing team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param team body model.Team true "Updated team data"
// @Success 200 {object} model.APIResponse "team updated"
// @Failure 400 {object} model.ErrorResponse "invalid data"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /team/{id} [put]
func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var team model.Team

	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Format data salah bre!",
			"data":    nil,
		})
	}

	if err := repo.UpdateTeam(c.Context(), id, team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Team berhasil di-update bre!",
	})
}

// DeleteTeam godoc
// @Summary Delete a team
// @Description Delete a team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} model.APIResponse "team deleted"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /team/{id} [delete]
func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := repo.DeleteTeam(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Team berhasil dihapus bre!",
	})
}
