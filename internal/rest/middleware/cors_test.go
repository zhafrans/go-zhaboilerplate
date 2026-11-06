package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"go-zhaboilerplate/internal/rest/middleware"
)

func TestCORS(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.CORS)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}
