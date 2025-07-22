package handler

import (
	"embeck/model"
	"embeck/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsers godoc
// @Summary Get All Users
// @Description Mendapatkan daftar semua user (Admin only)
// @Tags Users Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UsersListResponse "Daftar semua user berhasil diambil"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Token tidak valid"
// @Failure 403 {object} map[string]interface{} "Forbidden - Hanya admin yang dapat mengakses"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	// Get all users from repository
	users, err := repository.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.UsersListResponse{
		Message: "Users retrieved successfully",
		Data:    users,
		Total:   len(users),
	})
}

// GetUserByID godoc
// @Summary Get User by ID
// @Description Mendapatkan detail user berdasarkan ID (Admin only)
// @Tags Users Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example("64f123abc456def789012345")
// @Success 200 {object} model.UserProfile "User detail berhasil diambil"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Token tidak valid"
// @Failure 403 {object} map[string]interface{} "Forbidden - Hanya admin yang dapat mengakses"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get user from repository
	user, err := repository.GetUserByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user ID format") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return user profile (without password)
	profile := model.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// UpdateUser godoc
// @Summary Update User
// @Description Mengupdate data user (Admin only)
// @Tags Users Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example("64f123abc456def789012345")
// @Param request body model.UpdateUserRequest true "Data user yang akan diupdate"
// @Success 200 {object} model.UserProfile "User berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "Request data tidak valid"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Token tidak valid"
// @Failure 403 {object} map[string]interface{} "Forbidden - Hanya admin yang dapat mengakses"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "Username atau email sudah digunakan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Validation
	if err := validateUpdateUserRequest(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if user exists first
	existingUser, err := repository.GetUserByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user ID format") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	if existingUser == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Prepare update data - only update fields that are provided
	updateData := bson.M{}

	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Email != "" {
		updateData["email"] = strings.ToLower(req.Email)
	}
	if req.Role != "" {
		updateData["role"] = req.Role
	}

	// Update user in database
	_, err = repository.UpdateUser(c.Context(), id, updateData)
	if err != nil {
		if strings.Contains(err.Error(), "sudah digunakan") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "tidak ada data yang diupdate") || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or no changes made",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Get updated user data
	updatedUser, err := repository.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get updated user data",
		})
	}

	// Return updated user profile
	profile := model.UserProfile{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Menghapus user dari sistem (Admin only)
// @Tags Users Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example("64f123abc456def789012345")
// @Success 200 {object} map[string]interface{} "User berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Token tidak valid"
// @Failure 403 {object} map[string]interface{} "Forbidden - Hanya admin yang dapat mengakses"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Check if user exists first
	existingUser, err := repository.GetUserByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user ID format") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	if existingUser == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Delete user from database
	_, err = repository.DeleteUser(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ada data yang dihapus") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"user_id": id,
	})
}

// validateUpdateUserRequest validates update user request data
func validateUpdateUserRequest(req *model.UpdateUserRequest) error {
	// At least one field should be provided
	if req.Username == "" && req.Email == "" && req.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "At least one field (username, email, or role) must be provided")
	}

	// Username validation (if provided)
	if req.Username != "" {
		if len(req.Username) < 3 || len(req.Username) > 50 {
			return fiber.NewError(fiber.StatusBadRequest, "Username must be between 3 and 50 characters")
		}
		if strings.Contains(req.Username, " ") {
			return fiber.NewError(fiber.StatusBadRequest, "Username should not contain spaces")
		}
	}

	// Email validation (if provided)
	if req.Email != "" {
		// Simple email validation - can be enhanced with regex
		if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
		}
	}

	// Role validation (if provided)
	if req.Role != "" {
		if req.Role != "user" && req.Role != "admin" {
			return fiber.NewError(fiber.StatusBadRequest, "Role must be either 'user' or 'admin'")
		}
	}

	return nil
}
