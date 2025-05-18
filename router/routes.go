package router

import (
	"EMBECK/handler"

	"github.com/gofiber/fiber/v2"
)
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	//homepage
	api.Get("/", handler.Homepage)
	//endpoint team
	api.Get("/team", handler.GetAllTeams)
	api.Get("/team/:id", handler.GetTeamByID)
	app.Post("/api/team", handler.CreateTeam)
	app.Put("/api/team/:id", handler.UpdateTeam)
	app.Delete("/api/team/:id", handler.DeleteTeam)

	//endpoint player
	api.Get("/player", handler.GetAllPlayers)
	api.Get("/player/:id", handler.GetPlayerByID)
	app.Post("/api/player", handler.CreatePlayer)
	app.Put("/api/player/:id", handler.UpdatePlayer)
	app.Delete("/api/player/:id", handler.DeletePlayer)
}