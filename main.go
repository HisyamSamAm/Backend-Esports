package main

import (
	"EMBECK/config"
	"EMBECK/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	config.DB = config.MongoConnect(config.DBName)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	app := fiber.New()

app.Use(cors.New(cors.Config{
	AllowOrigins:     strings.Join(config.GetAllowedOrigins(), ","),
	AllowMethods:     "GET,POST,PUT,DELETE",
	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	AllowCredentials: true,
}))

router.SetupRoutes(app)

app.Use(func(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  fiber.StatusNotFound,
		"message": "Endpoint not found",
	})
})

port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}
log.Printf("Server running on port %s", port)
if err := app.Listen(":" + port); err != nil {
	log.Fatalf("Failed to start server: %v", err)
}

}