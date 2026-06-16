package middleware

import (
	"time"

	"go-user-api/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDMiddleware adds a unique request ID to each request
// This helps track requests across logs
func RequestIDMiddleware(c *fiber.Ctx) error {
	// Generate unique request ID
	requestID := uuid.New().String()

	// Add to response header
	c.Set("X-Request-ID", requestID)

	// Store in context for later use in handlers
	c.Locals("requestID", requestID)

	// Continue to next handler
	return c.Next()
}

// LoggingMiddleware logs details about each request
func LoggingMiddleware(c *fiber.Ctx) error {
	// Record start time
	start := time.Now()

	// Get request ID from context
	requestID := c.Locals("requestID").(string)

	// Log incoming request
	logger.Info("Incoming request",
		zap.String("requestID", requestID),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()),
	)

	// Call next handler
	err := c.Next()

	// Calculate request duration
	duration := time.Since(start).Milliseconds()

	// Log response details
	logger.Info("Request completed",
		zap.String("requestID", requestID),
		zap.Int("statusCode", c.Response().StatusCode()),
		zap.Int64("duration_ms", duration),
	)

	return err
}

// CORSMiddleware enables CORS (allows frontend to call backend)
func CORSMiddleware(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if c.Method() == "OPTIONS" {
		return c.SendStatus(200)
	}

	return c.Next()
}
