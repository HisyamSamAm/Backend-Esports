package router

import (
	"embeck/config/middleware"
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

	// ==================
	// Public Routes
	// ==================
	public := api.Group("/")
	public.Post("/auth/register", handler.Register)
	public.Post("/auth/login", handler.Login)
	public.Get("/tournaments", handler.GetAllTournamentsPublic)
	public.Get("/tournaments/:id", handler.GetTournamentWithDetailsByID)

	// ==================
	// Authenticated User Routes (User & Admin)
	// ==================
	authRequired := api.Group("/")
	authRequired.Use(middleware.AuthMiddleware())
	authRequired.Get("/auth/profile", handler.GetProfile) // Now requires auth
	authRequired.Post("/tickets/purchase", handler.HandlePurchaseTicket)
	authRequired.Get("/me/tickets", handler.HandleGetUserTickets)

	// ==================
	// Admin Only Routes
	// ==================
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// Player Management (Admin)
	admin.Get("/players", handler.GetAllPlayers)
	admin.Post("/players", handler.CreatePlayer)
	admin.Get("/players/:id", handler.GetPlayerByID)
	admin.Put("/players/:id", handler.UpdatePlayer)
	admin.Delete("/players/:id", handler.DeletePlayer)

	// Team Management (Admin)
	admin.Get("/teams", handler.GetAllTeams)
	admin.Post("/teams", handler.CreateTeam)
	admin.Get("/teams/:id", handler.GetTeamByID)
	admin.Put("/teams/:id", handler.UpdateTeam)
	admin.Delete("/teams/:id", handler.DeleteTeam)

	// Tournament Management (Admin)
	admin.Get("/tournaments", handler.GetAllTournaments)
	admin.Post("/tournaments", handler.CreateTournament)
	admin.Get("/tournaments/:id", handler.GetTournamentByID)
	admin.Put("/tournaments/:id", handler.UpdateTournament)
	admin.Delete("/tournaments/:id", handler.DeleteTournament)

	// Match Management (Admin)
	admin.Get("/matches", handler.GetAllMatches)
	admin.Post("/matches", handler.CreateMatch)
	admin.Get("/matches/:id", handler.GetMatchByID)
	admin.Put("/matches/:id", handler.UpdateMatch)
	admin.Delete("/matches/:id", handler.DeleteMatch)

	// User Management (Admin)
	admin.Get("/users", handler.GetAllUsers)
	admin.Get("/users/:id", handler.GetUserByID)
	admin.Put("/users/:id", handler.UpdateUser)
	admin.Delete("/users/:id", handler.DeleteUser)

	// Upload routes (Admin)
	admin.Post("/upload/team-logo", handler.UploadTeamLogo)
	admin.Post("/upload/player-avatar", handler.UploadPlayerAvatar)
}
