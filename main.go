package main

import (
	"EMBECK/config"
	_ "EMBECK/docs"
	"EMBECK/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title EMBECK API
// @version 1.0
// @description API untuk manajemen esports
// @host localhost:1010
// @BasePath /api

func init() {
	_ = godotenv.Load()

	config.ConnectDB() // Initialize MongoDB connection
}

func main() {
	app := fiber.New()

	app.Use(logger.New())
	// CORS setup
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.GetAllowedOrigins(), ","),
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Setup routes
	router.SetupRoutes(app)

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Fallback route for 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Endpoint not found",
		})
	})

	// Ambil port dari .env, default ke 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
