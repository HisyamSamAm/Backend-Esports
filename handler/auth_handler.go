package handler

import (
	"embeck/model"
	"embeck/pkg/password"
	"embeck/repository"
	"regexp"
	"strings"
	"time"

	"embeck/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register godoc
// @Summary Register New User
// @Description Mendaftarkan user baru ke dalam sistem
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Data registrasi user"
// @Success 201 {object} model.UserResponse "User berhasil didaftarkan"
// @Failure 400 {object} map[string]interface{} "Request data tidak valid"
// @Failure 409 {object} map[string]interface{} "Email atau username sudah digunakan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/register [post]
func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Validation
	if err := validateRegisterRequest(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	// Create user model
	user := model.User{
		Username: req.Username,
		Email:    strings.ToLower(req.Email),
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	// Create user in database
	insertedID, err := repository.CreateUser(c.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "sudah terdaftar") || strings.Contains(err.Error(), "sudah digunakan") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Convert ObjectID to string
	var userID string
	if oid, ok := insertedID.(primitive.ObjectID); ok {
		userID = oid.Hex()
	}

	return c.Status(fiber.StatusCreated).JSON(model.UserResponse{
		Message: "User registered successfully",
		UserID:  userID,
	})
}

// Login godoc
// @Summary User Login
// @Description Login user dan mendapatkan PASETO token untuk autentikasi
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Data login user"
// @Success 200 {object} model.AuthResponse "Login berhasil dengan token"
// @Failure 400 {object} map[string]interface{} "Request data tidak valid"
// @Failure 401 {object} map[string]interface{} "Kredensial tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/login [post]
func Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	user, err := repository.GetUserByEmail(c.Context(), strings.ToLower(req.Email))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := password.CheckPassword(user.Password, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// 🔒 Generate PASETO token (valid for 24 hours as per auth pkg)
	token, err := auth.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.AuthResponse{
		Message:  "Login successful",
		Token:    token,
		Role:     user.Role,
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	})
}

// GetProfile godoc
// @Summary Get User Profile (Demo Mode)
// @Description Mendapatkan profil user untuk testing - tidak memerlukan authentication (demo mode)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user_id query string false "User ID untuk mendapatkan profil (opsional)"
// @Success 200 {object} model.UserProfile "Profil user"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/profile [get]
func GetProfile(c *fiber.Ctx) error {
	// Get user ID from query parameter (since no auth required)
	userID := c.Query("user_id")

	// If no user_id provided, return default admin user info
	if userID == "" {
		// Return default admin profile for demo
		demoObjectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011") // Demo ObjectID
		profile := model.UserProfile{
			ID:        demoObjectID,
			Username:  "demo-admin",
			Email:     "admin@embeck.com",
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return c.Status(fiber.StatusOK).JSON(profile)
	}

	// Get user from database if ID provided
	user, err := repository.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user profile",
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

// validateRegisterRequest validates registration request data
func validateRegisterRequest(req *model.RegisterRequest) error {
	// Username validation
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Username must be between 3 and 50 characters")
	}

	// Username should not contain spaces
	if strings.Contains(req.Username, " ") {
		return fiber.NewError(fiber.StatusBadRequest, "Username should not contain spaces")
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}

	// Password validation
	if len(req.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 6 characters")
	}

	return nil
}
