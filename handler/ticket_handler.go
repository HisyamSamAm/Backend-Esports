package handler

import (
	"embeck/model"
	"embeck/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTicket godoc
// @Summary Create New Ticket Type
// @Description Membuat jenis tiket baru untuk turnamen. Setiap tournament bisa memiliki berbagai jenis tiket dengan harga dan kapasitas yang berbeda.
// @Tags Tickets
// @Accept json
// @Produce json
// @Param request body model.TicketRequest true "Data tiket yang akan dibuat"
// @Success 201 {object} model.TicketResponse "Ticket berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Request data tidak valid"
// @Failure 409 {object} map[string]interface{} "Tournament tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/tickets [post]
func CreateTicket(c *fiber.Ctx) error {
	var req model.TicketRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Validation
	if req.TournamentID == "" || req.Price < 0 || req.QuantityAvailable < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "tournament_id is required, price and quantity_available must be non-negative",
		})
	}

	// Convert string ID to ObjectID
	tournamentObjID, err := primitive.ObjectIDFromHex(req.TournamentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament_id format",
		})
	}

	// Create ticket model
	ticket := model.Ticket{
		TournamentID:      tournamentObjID,
		Price:             req.Price,
		QuantityAvailable: req.QuantityAvailable,
		Description:       req.Description,
	}

	insertedID, err := repository.CreateTicket(c.Context(), ticket)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal menambahkan ticket: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.TicketResponse{
		Message:  "Ticket type created successfully",
		TicketID: insertedID.(primitive.ObjectID).Hex(),
	})
}

// GetAllTickets godoc
// @Summary Get All Ticket Types
// @Description Mendapatkan daftar semua jenis tiket dengan opsi filter dan populate tournament details
// @Tags Tickets
// @Accept json
// @Produce json
// @Param tournament_id query string false "Filter berdasarkan Tournament ID"
// @Param populate query bool false "Include tournament details (true/false)" default(false)
// @Success 200 {array} model.Ticket "List semua ticket types (tanpa populate)"
// @Success 200 {array} model.TicketWithTournament "List semua ticket types dengan tournament details (dengan populate=true)"
// @Failure 400 {object} map[string]interface{} "Tournament ID format tidak valid"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data dari database"
// @Router /api/admin/tickets [get]
func GetAllTickets(c *fiber.Ctx) error {
	// Get optional tournament_id filter from query params
	tournamentID := c.Query("tournament_id")
	populate := c.QueryBool("populate", false)

	if populate {
		tickets, err := repository.GetAllTicketsWithTournament(c.Context(), tournamentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal mengambil data tickets dengan detail tournament dari database",
			})
		}
		return c.Status(fiber.StatusOK).JSON(tickets)
	}

	tickets, err := repository.GetAllTickets(c.Context(), tournamentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data tickets dari database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tickets)
}

// GetTicketByID godoc
// @Summary Get Ticket Type By ID
// @Description Mendapatkan detail jenis tiket berdasarkan ID dengan opsi populate tournament details
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID (ObjectID format)" example("64f123abc456def789012345")
// @Param populate query bool false "Include tournament details (true/false)" default(false)
// @Success 200 {object} model.Ticket "Detail ticket type (tanpa populate)"
// @Success 200 {object} model.TicketWithTournament "Detail ticket type dengan tournament details (dengan populate=true)"
// @Failure 400 {object} map[string]interface{} "Ticket ID format tidak valid"
// @Failure 404 {object} map[string]interface{} "Ticket type tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/tickets/{id} [get]
func GetTicketByID(c *fiber.Ctx) error {
	id := c.Params("id")
	populate := c.QueryBool("populate", false)

	if populate {
		ticket, err := repository.GetTicketWithTournamentByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if ticket == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Ticket type not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(ticket)
	}

	ticket, err := repository.GetTicketByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if ticket == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket type not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ticket)
}

// UpdateTicket godoc
// @Summary Update Ticket Type
// @Description Memperbarui detail jenis tiket. Field yang tidak diisi akan tetap menggunakan nilai lama.
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID (ObjectID format)" example("64f123abc456def789012345")
// @Param request body model.TicketRequest true "Data tiket yang akan diupdate (partial update supported)"
// @Success 200 {object} model.TicketResponse "Ticket berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "Request data atau ID format tidak valid"
// @Failure 404 {object} map[string]interface{} "Ticket type tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/tickets/{id} [put]
func UpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.TicketRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Create update model (only update non-empty/non-zero fields)
	update := model.Ticket{}

	if req.TournamentID != "" {
		tournamentObjID, err := primitive.ObjectIDFromHex(req.TournamentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid tournament_id format",
			})
		}
		update.TournamentID = tournamentObjID
	}

	if req.Price >= 0 {
		update.Price = req.Price
	}

	if req.QuantityAvailable >= 0 {
		update.QuantityAvailable = req.QuantityAvailable
	}

	if req.Description != "" {
		update.Description = req.Description
	}

	_, err := repository.UpdateTicket(c.Context(), id, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Error updating ticket %s: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.TicketResponse{
		Message: "Ticket type updated successfully",
	})
}

// DeleteTicket godoc
// @Summary Delete Ticket Type
// @Description Menghapus jenis tiket dari database secara permanen. Operasi ini tidak dapat dibatalkan.
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID (ObjectID format)" example("64f123abc456def789012345")
// @Success 200 {object} model.TicketResponse "Ticket berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "Ticket ID format tidak valid"
// @Failure 404 {object} map[string]interface{} "Ticket type tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/tickets/{id} [delete]
func DeleteTicket(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeleteTicket(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Ticket dengan ID %s tidak ditemukan: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.TicketResponse{
		Message: "Ticket type deleted successfully",
	})
}

// GetTicketsByTournament godoc
// @Summary Get Tickets by Tournament ID
// @Description Mendapatkan semua jenis tiket yang tersedia untuk tournament tertentu
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "Tournament ID"
// @Success 200 {array} model.Ticket "List tiket untuk tournament"
// @Failure 400 {object} map[string]interface{} "ID tournament tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/tournaments/{id}/tickets [get]
func GetTicketsByTournament(c *fiber.Ctx) error {
	tournamentIDStr := c.Params("id")

	tournamentID, err := primitive.ObjectIDFromHex(tournamentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	tickets, err := repository.GetTicketsByTournament(tournamentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get tickets: %v", err),
		})
	}

	return c.JSON(tickets)
}

// GetTicketsByUser godoc
// @Summary Get Tickets by User ID
// @Description Mendapatkan semua tiket yang dimiliki oleh user tertentu
// @Tags Tickets
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} model.Ticket "List tiket user"
// @Failure 400 {object} map[string]interface{} "ID user tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/tickets/user/{userId} [get]
func GetTicketsByUser(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	tickets, err := repository.GetTicketsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get user tickets: %v", err),
		})
	}

	return c.JSON(tickets)
}
