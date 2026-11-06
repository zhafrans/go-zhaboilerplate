package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/rest"
	"go-zhaboilerplate/internal/rest/middleware"
	"go-zhaboilerplate/internal/rest/mocks"
)

func TestAuthHandlerRegister(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.AuthService)
	mockUser := domain.User{ID: 1, Email: "user@example.com", IsActive: true}

	mockSvc.On("Register", mock.Anything, "user@example.com", "password123").Return(mockUser, nil)

	rest.NewAuthHandler(app, mockSvc, nil)

	payload := rest.RegisterRequest{
		Email:    "user@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestAuthHandlerLogin(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.AuthService)
	mockResp := domain.LoginResponse{
		Token: "jwt-token-string",
		User:  domain.User{ID: 1, Email: "user@example.com", IsActive: true},
	}

	mockSvc.On("Login", mock.Anything, "user@example.com", "password123").Return(mockResp, nil)

	rest.NewAuthHandler(app, mockSvc, nil)

	payload := rest.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandlerMe(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.AuthService)
	mockUser := domain.User{ID: 1, Email: "user@example.com", IsActive: true}

	mockSvc.On("GetProfile", mock.Anything, int64(1)).Return(mockUser, nil)

	testAuthMiddleware := func(c *fiber.Ctx) error {
		c.Locals(middleware.ContextUserIDKey, int64(1))
		c.Locals(middleware.ContextEmailKey, "user@example.com")
		return c.Next()
	}

	rest.NewAuthHandler(app, mockSvc, testAuthMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
