package routes

import (
	"go-user-api/internal/handler"
	"go-user-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	// Apply global middleware
	app.Use(middleware.RequestIDMiddleware)
	app.Use(middleware.LoggingMiddleware)
	app.Use(middleware.CORSMiddleware)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	// User routes
	api := app.Group("/users")

	// POST /users - Create user
	api.Post("", userHandler.CreateUser)

	// GET /users - Get all users
	api.Get("", userHandler.GetAllUsers)

	// GET /users/:id - Get user by ID
	api.Get("/:id", userHandler.GetUser)

	// PUT /users/:id - Update user
	api.Put("/:id", userHandler.UpdateUser)

	// DELETE /users/:id - Delete user
	api.Delete("/:id", userHandler.DeleteUser)
}
