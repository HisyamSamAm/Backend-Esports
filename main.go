package main

import (
	"embeck/config"
	"embeck/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	_ "embeck/docs"
)

func init() {
	// Load .env file if exists
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
}

// @title ESports Management API
// @version 1.0
// @description This is the API for the ESports Management platform.
// @contact.name API Support
// @contact.email support@embeck.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:1010
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a PASETO token.
func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "EMBECK API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Test database connection
	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Setup Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.GetAllowedOrigins(), ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Setup routes
	router.SetupRoutes(app)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":  "Hello, EMBECK Backend is running! 🚀",
			"version":  "1.0.0",
			"status":   "active",
			"database": "connected",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"uptime": "running",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
