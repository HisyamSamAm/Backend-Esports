package handler

import (
	"EMBECK/auth"
	repo "EMBECK/repository"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	users, err := repo.GetAllUsers(c.Context())
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
		"data":    users,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := repo.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "user gak ketemu bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success ambil 1 user!",
		"data":    user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := repo.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data user bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil nambahin data user bre!",
		"data":    user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedData model.User

	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := repo.UpdateUser(c.Context(), id, updatedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update data user bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil update data user bre!",
	})
}
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repo.DeleteUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus data user bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil hapus data user bre!",
	})
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with username, email, password and role
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User registration data"
// @Success 200 {object} model.UserResponse "user registered successfully"
// @Failure 400 {object} model.ErrorResponse "invalid data or admin role not allowed"
// @Failure 500 {object} model.ErrorResponse "server error"
// @Router /register [post]
func RegisterHandler(c *fiber.Ctx) error {
	return auth.Register(c)
}

// LoginHandler godoc
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse "login successful"
// @Failure 400 {object} model.ErrorResponse "invalid data"
// @Failure 401 {object} model.ErrorResponse "invalid credentials"
// @Failure 403 {object} model.ErrorResponse "invalid user role"
// @Router /login [post]
func LoginHandler(c *fiber.Ctx) error {
	return auth.Login(c)
}
