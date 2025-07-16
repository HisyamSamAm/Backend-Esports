package router

import (
	"embeck/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	// Swagger documentation route
	app.Get("/docs/*", swagger.HandlerDefault)

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

	// API group
	api := app.Group("/api")

	// Authentication routes (public - optional untuk testing)
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Get("/profile", handler.GetProfile) // Removed auth middleware

	// Upload routes (now public)
	upload := api.Group("/upload")
	upload.Post("/team-logo", handler.UploadTeamLogo)
	upload.Post("/player-avatar", handler.UploadPlayerAvatar)

	// Player routes (now public - no auth required)
	api.Get("/players", handler.GetAllPlayers)
	api.Get("/players/:id", handler.GetPlayerByID)
	api.Post("/players", handler.CreatePlayer)
	api.Put("/players/:id", handler.UpdatePlayer)
	api.Delete("/players/:id", handler.DeletePlayer)

	// Team routes (now public - no auth required)
	api.Get("/teams", handler.GetAllTeams)
	api.Get("/teams/:id", handler.GetTeamByID)
	api.Post("/teams", handler.CreateTeam)
	api.Put("/teams/:id", handler.UpdateTeam)
	api.Delete("/teams/:id", handler.DeleteTeam)

	// Tournament routes (now public - no auth required)
	api.Get("/tournaments", handler.GetAllTournamentsPublic)
	api.Get("/tournaments/:id", handler.GetTournamentWithDetailsByID)
	api.Post("/tournaments", handler.CreateTournament)
	api.Put("/tournaments/:id", handler.UpdateTournament)
	api.Delete("/tournaments/:id", handler.DeleteTournament)

	// Match routes (now public - no auth required)
	api.Get("/matches", handler.GetAllMatches)
	api.Get("/matches/:id", handler.GetMatchByID)
	api.Post("/matches", handler.CreateMatch)
	api.Put("/matches/:id", handler.UpdateMatch)
	api.Delete("/matches/:id", handler.DeleteMatch)

	// Ticket routes (now public - no auth required)
	api.Get("/tickets", handler.GetAllTickets)
	api.Get("/tickets/:id", handler.GetTicketByID)
	api.Post("/tickets", handler.CreateTicket)
	api.Put("/tickets/:id", handler.UpdateTicket)
	api.Delete("/tickets/:id", handler.DeleteTicket)
	api.Get("/tickets/tournament/:tournamentId", handler.GetTicketsByTournament)
	api.Get("/tickets/user/:userId", handler.GetTicketsByUser)

	// Legacy admin routes (redirected to main routes for compatibility)
	admin := api.Group("/admin")
	admin.Get("/players", handler.GetAllPlayers)
	admin.Get("/players/:id", handler.GetPlayerByID)
	admin.Post("/players", handler.CreatePlayer)
	admin.Put("/players/:id", handler.UpdatePlayer)
	admin.Delete("/players/:id", handler.DeletePlayer)

	admin.Get("/teams", handler.GetAllTeams)
	admin.Get("/teams/:id", handler.GetTeamByID)
	admin.Post("/teams", handler.CreateTeam)
	admin.Put("/teams/:id", handler.UpdateTeam)
	admin.Delete("/teams/:id", handler.DeleteTeam)

	admin.Get("/tournaments", handler.GetAllTournaments)
	admin.Get("/tournaments/:id", handler.GetTournamentByID)
	admin.Post("/tournaments", handler.CreateTournament)
	admin.Put("/tournaments/:id", handler.UpdateTournament)
	admin.Delete("/tournaments/:id", handler.DeleteTournament)

	admin.Get("/matches", handler.GetAllMatches)
	admin.Get("/matches/:id", handler.GetMatchByID)
	admin.Post("/matches", handler.CreateMatch)
	admin.Put("/matches/:id", handler.UpdateMatch)
	admin.Delete("/matches/:id", handler.DeleteMatch)

	admin.Get("/tickets", handler.GetAllTickets)
	admin.Get("/tickets/:id", handler.GetTicketByID)
	admin.Post("/tickets", handler.CreateTicket)
	admin.Put("/tickets/:id", handler.UpdateTicket)
	admin.Delete("/tickets/:id", handler.DeleteTicket)
}
