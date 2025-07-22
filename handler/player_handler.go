package handler

import (
	"embeck/model"
	"embeck/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetAllPlayers godoc
// @Summary Get All Players
// @Description Mendapatkan daftar semua pemain Mobile Legends
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Player "Daftar semua pemain"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/players [get]
func GetAllPlayers(c *fiber.Ctx) error {
	players, err := repository.GetAllPlayers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data players dari database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(players)
}

// GetPlayerByID godoc
// @Summary Get Player By ID
// @Description Mendapatkan detail pemain berdasarkan ID
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Player ID" example("64f123abc456def789012345")
// @Success 200 {object} model.Player "Detail pemain"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Player tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/players/{id} [get]
func GetPlayerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	player, err := repository.GetPlayerByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if player == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Player not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(player)
}

// CreatePlayer godoc
// @Summary Create New Player
// @Description Menambahkan pemain baru ke database
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.PlayerRequest true "Player data"
// @Success 201 {object} model.PlayerResponse "Player berhasil dibuat"
// @Failure 400 {object} model.ErrorResponse "Request data tidak valid"
// @Failure 409 {object} model.ErrorResponse "Player sudah ada"
// @Router /api/admin/players [post]
func CreatePlayer(c *fiber.Ctx) error {
	var req model.PlayerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request data",
		})
	}

	// Validation
	if req.Name == "" || req.MLNickname == "" || req.MLID == "" || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "missing_fields",
			Message: "All fields (name, ml_nickname, ml_id, status) are required",
		})
	}

	// Create player model
	player := model.Player{
		Name:       req.Name,
		MLNickname: req.MLNickname,
		MLID:       req.MLID,
		Status:     req.Status,
	}

	insertedID, err := repository.InsertPlayer(c.Context(), player)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
			Error:   "db_conflict",
			Message: fmt.Sprintf("Gagal menambahkan player: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Player created successfully",
		"player_id": insertedID,
	})
}

// UpdatePlayer godoc
// @Summary Update Player
// @Description Memperbarui detail pemain
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Player ID" example("64f123abc456def789012345")
// @Param request body model.PlayerRequest true "Player data"
// @Success 200 {object} model.PlayerResponse "Player berhasil diupdate"
// @Failure 400 {object} model.ErrorResponse "Request data tidak valid"
// @Failure 404 {object} model.ErrorResponse "Player tidak ditemukan"
// @Router /api/admin/players/{id} [put]
func UpdatePlayer(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.PlayerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request data",
		})
	}

	// Create update model (only update non-empty fields)
	update := model.Player{}
	if req.Name != "" {
		update.Name = req.Name
	}
	if req.MLNickname != "" {
		update.MLNickname = req.MLNickname
	}
	if req.MLID != "" {
		update.MLID = req.MLID
	}
	if req.Status != "" {
		update.Status = req.Status
	}

	_, err := repository.UpdatePlayer(c.Context(), id, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "update_failed",
			Message: fmt.Sprintf("Error updating player %s: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player updated successfully",
	})
}

// DeletePlayer godoc
// @Summary Delete Player
// @Description Menghapus pemain dari database
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Player ID" example("64f123abc456def789012345")
// @Success 200 {object} map[string]interface{} "Player berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Player tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/players/{id} [delete]
func DeletePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeletePlayer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Player dengan ID %s tidak ditemukan: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player deleted successfully",
	})
}
