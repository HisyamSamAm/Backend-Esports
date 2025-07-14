package auth

import (
	"EMBECK/model"
	usr "EMBECK/repository"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	// Authenticate user against database or auth system
	dbUser, err := usr.AuthenticateUser(c.Context(), user.Username, user.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid username or password",
			"data":    nil,
		})
	}

	// Check user role
	if dbUser.Role == "admin" {
		// Redirect to admin dashboard
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Berhasil login sebagai admin",
			"data":    dbUser,
		})
	} else if dbUser.Role == "user" {
		// Redirect to user dashboard
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Berhasil login sebagai user",
			"data":    dbUser,
		})
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Invalid user role",
			"data":    nil,
		})
	}
}