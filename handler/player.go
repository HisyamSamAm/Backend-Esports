package handler

import (
	"EMBECK/model"
	repo "EMBECK/repository"

	"github.com/gofiber/fiber/v2"
)

// GetAllPlayers godoc
// @Summary Get all players
// @Description Get list of all players
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {object} model.PlayersResponse "success"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /player [get]
func GetAllPlayers(c *fiber.Ctx) error {
	players, err := repo.GetAllPlayers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal ambil data players!",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil ambil data players!",
		"data":    players,
	})
}

// GetPlayerByID godoc
// @Summary Get player by ID
// @Description Get a single player by its ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} model.PlayerResponse "success"
// @Failure 404 {object} model.ErrorResponse "player not found"
// @Router /player/{id} [get]
func GetPlayerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	player, err := repo.GetPlayerByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Player tidak ditemukan!",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil ambil player!",
		"data":    player,
	})
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Create a new player with the provided data
// @Tags players
// @Accept json
// @Produce json
// @Param player body model.Player true "Player data"
// @Success 201 {object} model.APIResponse "player created"
// @Failure 400 {object} model.ErrorResponse "invalid data"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /player [post]
func CreatePlayer(c *fiber.Ctx) error {
	var player model.Player
	if err := c.BodyParser(&player); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse data!",
			"error":   err.Error(),
		})
	}

	err := repo.CreatePlayer(c.Context(), player)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal simpan player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Player berhasil ditambahkan!",
		"data":    player,
	})
}

// UpdatePlayer godoc
// @Summary Update a player
// @Description Update an existing player by ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param player body model.Player true "Updated player data"
// @Success 200 {object} model.APIResponse "player updated"
// @Failure 400 {object} model.ErrorResponse "invalid data"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /player/{id} [put]
func UpdatePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	var updatedPlayer model.Player
	if err := c.BodyParser(&updatedPlayer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse data!",
			"error":   err.Error(),
		})
	}

	// Pastikan ID tidak berubah ke kosong/null
	updatedPlayer.ID = id

	err := repo.UpdatePlayer(c.Context(), id, updatedPlayer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player berhasil diupdate!",
		"data":    updatedPlayer,
	})
}

// DeletePlayer godoc
// @Summary Delete a player
// @Description Delete a player by ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} model.APIResponse "player deleted"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /player/{id} [delete]
func DeletePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repo.DeletePlayer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player berhasil dihapus!",
	})
}
