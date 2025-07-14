package auth

import (
	"EMBECK/model"
	pwd "EMBECK/pkg"
	repo "EMBECK/repository"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest

	// Parsing JSON ke struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validasi kosong
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	// Hash password
	hashedPassword, err := pwd.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	userLogin := model.UserLogin{
    Username: req.Username,
    Password: hashedPassword,
    Role:     "user",
}

// Insert ke database
	id, err := repo.InsertUser(c.Context(), userLogin)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"id":      id,
	})
}

