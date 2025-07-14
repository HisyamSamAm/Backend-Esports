package auth

import (
	"EMBECK/config/middleware"
	"EMBECK/model"
	pwd "EMBECK/pkg"
	repo "EMBECK/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var req model.UserLogin

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := repo.FindUserByUsername(c.Context(), req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username not found"})
	}

	// Cek password input hash yang tersimpan
	if !pwd.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong password"})
	}
<<<<<<< HEAD
}
=======

	// Generate token PASETO
	token, err := middleware.EncodeWithRoleHours(user.Role, user.Username, 2)
	if err != nil {
		fmt.Println("Token generation error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   token,
	})
}
>>>>>>> 4d86ebc93116bd04a770e6a9be3c2754399fe534
