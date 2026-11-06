package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"go-zhaboilerplate/internal/rest/middleware"
	"go-zhaboilerplate/internal/token"
)

func TestJWTMiddleware(t *testing.T) {
	secret := "jwt-test-secret"
	mw := middleware.JWTMiddleware(secret)

	app := fiber.New()
	app.Get("/protected", mw, func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserID(c)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), userID)

		email, err := middleware.GetUserEmail(c)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", email)

		return c.SendString("ok")
	})

	t.Run("missing header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("malformed header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Token abcxyz")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("valid token", func(t *testing.T) {
		tok, err := token.GenerateToken(123, "user@example.com", secret, time.Hour)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
