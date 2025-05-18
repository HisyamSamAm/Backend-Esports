package handler

import (
	"EMBECK/controller"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllTickets(c *fiber.Ctx) error {
	tickets, err := controller.GetAllTickets(c.Context())
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
		"data":    tickets,
	})
}
func GetTicketByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ticket, err := controller.GetTicketByID(c.Context(), id)
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
		"data":    ticket,
	})
}

func CreateTicket(c *fiber.Ctx) error {
	var ticket model.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}
	err := controller.CreateTicket(c.Context(), ticket)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data ticket bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "success nambahin data ticket bre!",
	})
}

func UpdateTicket(c *fiber.Ctx) error {
	var ticket model.Ticket
	id := c.Params("id")
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}
	err := controller.UpdateTicket(c.Context(), id, ticket)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data ticket bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success nambahin data ticket bre!",
	})
}

func DeleteTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.DeleteTicket(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data ticket bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success nambahin data ticket bre!",
	})
}