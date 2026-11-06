package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go-zhaboilerplate/internal/rest/response"
	"go-zhaboilerplate/internal/token"
)

const (
	ContextUserIDKey = "user_id"
	ContextEmailKey  = "email"
)

var (
	ErrUserContextNotFound = errors.New("user context not found")
)

// JWTMiddleware creates a fiber middleware that validates JWT bearer token
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return response.Error(c, http.StatusUnauthorized, "Missing Authorization header", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'", nil)
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := token.ValidateToken(tokenString, secret)
		if err != nil {
			return response.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
		}

		c.Locals(ContextUserIDKey, claims.UserID)
		c.Locals(ContextEmailKey, claims.Email)

		return c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from fiber context locals
func GetUserID(c *fiber.Ctx) (int64, error) {
	val := c.Locals(ContextUserIDKey)
	if val == nil {
		return 0, ErrUserContextNotFound
	}

	id, ok := val.(int64)
	if !ok {
		return 0, ErrUserContextNotFound
	}

	return id, nil
}

// GetUserEmail retrieves the authenticated user's email from fiber context locals
func GetUserEmail(c *fiber.Ctx) (string, error) {
	val := c.Locals(ContextEmailKey)
	if val == nil {
		return "", ErrUserContextNotFound
	}

	email, ok := val.(string)
	if !ok {
		return "", ErrUserContextNotFound
	}

	return email, nil
}
