package handler

import (
	"EMBECK/controller"

	"github.com/gofiber/fiber/v2"
)

func GetAllTeams(c *fiber.Ctx) error {
	team, err := controller.GetAllTeams(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
			"message": "error nih servernya bre!",
			"data":    nil,
		})
}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"message": "success ngambil data bre!",
		"data":    team,
	})
}	