package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/rest"
	"go-zhaboilerplate/internal/rest/mocks"
)

func TestUserHandlerFetch(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)
	mockUsers := []domain.User{
		{ID: 1, Email: "user1@example.com", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockSvc.On("Fetch", mock.Anything, "", int64(10)).Return(mockUsers, "nextCursor", nil)

	rest.NewUserHandler(app, mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandlerGetByID(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)
	mockUser := domain.User{ID: 1, Email: "user1@example.com", IsActive: true}

	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(mockUser, nil)

	rest.NewUserHandler(app, mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandlerStore(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)

	mockSvc.On("Store", mock.Anything, mock.Anything).Return(nil)

	rest.NewUserHandler(app, mockSvc, nil)

	payload := rest.CreateUserRequest{
		Email:    "new@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUserHandlerUpdate(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)

	mockSvc.On("Update", mock.Anything, mock.Anything).Return(nil)

	rest.NewUserHandler(app, mockSvc, nil)

	payload := rest.UpdateUserRequest{
		Email: "updated@example.com",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandlerDelete(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	rest.NewUserHandler(app, mockSvc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandlerSetActive(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mocks.UserService)

	mockSvc.On("SetActive", mock.Anything, int64(1), false).Return(nil)

	rest.NewUserHandler(app, mockSvc, nil)

	isActive := false
	payload := rest.SetActiveRequest{
		IsActive: &isActive,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/users/1/active", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
