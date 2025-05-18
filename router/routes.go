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

	//endpoint tournament
	api.Get("/tournament", handler.GetAllTournaments)
	api.Get("/tournament/:id", handler.GetTournamentByID)
	app.Post("/api/tournament", handler.CreateTournament)
	app.Put("/api/tournament/:id", handler.UpdateTournament)
	app.Delete("/api/tournament/:id", handler.DeleteTournament)

	//endpoint score match
	api.Get("/score-match", handler.GetAllScoreMatches)
	api.Get("/score-match/:id", handler.GetScoreMatchByID)
	app.Post("/api/score-match", handler.CreateScoreMatch)
	app.Put("/api/score-match/:id", handler.UpdateScoreMatch)
	app.Delete("/api/score-match/:id", handler.DeleteScoreMatch)
}