package handler

import "github.com/gofiber/fiber/v2"

func Homepage(c *fiber.Ctx) error {
	return c.SendString("Welcome to the EMBECK API")
}