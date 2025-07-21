package handler

import (
	"embeck/model"
	"embeck/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandlePurchaseTicket handles the logic for a user purchasing a ticket for a match.
func HandlePurchaseTicket(c *fiber.Ctx) error {
	var req model.UserTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_request", Message: "Cannot parse JSON"})
	}

	if req.MatchID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "missing_field", Message: "match_id is required"})
	}

	matchObjID, err := primitive.ObjectIDFromHex(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Error: "invalid_id", Message: "Invalid match_id format"})
	}

	// In a real application, UserID would come from the JWT token.
	// For this simplified version, we'll extract it, but acknowledge it's a placeholder.
	// We'll assume the middleware has validated the token and the user's role.
	claims, ok := c.Locals("claims").(*model.TokenClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Error: "unauthorized", Message: "Invalid or missing token claims"})
	}

	userObjID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Error: "internal_error", Message: "Could not parse user ID from token"})
	}

	// Call repository to perform the purchase
	newTicket, err := repository.PurchaseTicket(c.Context(), userObjID, matchObjID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Error: "not_found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "already purchased") {
			return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{Error: "conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Error: "database_error", Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(newTicket)
}

// HandleGetUserTickets retrieves all tickets for the currently authenticated user.
func HandleGetUserTickets(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*model.TokenClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Error: "unauthorized", Message: "Invalid or missing token claims"})
	}

	userObjID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Error: "internal_error", Message: "Could not parse user ID from token"})
	}

	tickets, err := repository.GetTicketsByUserID(c.Context(), userObjID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Error: "database_error", Message: err.Error()})
	}

	if len(tickets) == 0 {
		return c.Status(fiber.StatusOK).JSON([]model.UserTicketResponse{})
	}

	return c.Status(fiber.StatusOK).JSON(tickets)
}
