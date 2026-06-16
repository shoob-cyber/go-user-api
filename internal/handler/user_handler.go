package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go-user-api/internal/logger"
	"go-user-api/internal/models"
	"go-user-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service   *service.UserService
	validator *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateUser handles POST /users
// This is called when frontend sends new user data
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if err := h.validator.Struct(req); err != nil {
		logger.Warn("Validation failed", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed: " + err.Error(),
		})
	}

	// Call service to create user
	user, err := h.service.CreateUser(req.Name, req.DOB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	logger.Info("User creation successful", zap.Int32("id", user.ID))
	return c.Status(http.StatusCreated).JSON(user)
}

// GetUser handles GET /users/:id
// This is called when frontend requests a specific user
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Get ID from URL parameter
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid user ID", zap.String("id", idStr))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Call service to get user
	user, err := h.service.GetUser(int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(user)
}

// UpdateUser handles PUT /users/:id
// This is called when frontend updates user information
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if err := h.validator.Struct(req); err != nil {
		logger.Warn("Validation failed", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed: " + err.Error(),
		})
	}

	// Call service to update user
	user, err := h.service.UpdateUser(int32(id), req.Name, req.DOB)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(user)
}

// DeleteUser handles DELETE /users/:id
// This is called when frontend wants to delete a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Call service to delete user
	err = h.service.DeleteUser(int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	// Return 204 No Content
	return c.SendStatus(http.StatusNoContent)
}

// GetAllUsers handles GET /users
// This is called when frontend wants to see all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Call service to get all users
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.JSON(users)
}

