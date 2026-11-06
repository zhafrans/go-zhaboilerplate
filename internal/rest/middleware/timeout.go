package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetRequestContextWithTimeout will set the request context with timeout for incoming requests
func SetRequestContextWithTimeout(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel()

		c.SetUserContext(ctx)
		return c.Next()
	}
}
