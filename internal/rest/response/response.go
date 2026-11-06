package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"go-zhaboilerplate/domain"
)

// Response represents a standard JSON API response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// Meta represents pagination or cursor metadata
type Meta struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int64  `json:"limit,omitempty"`
}

// Success sends a standardized success JSON response
func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a standardized success JSON response with metadata (e.g. cursor pagination)
func SuccessWithMeta(c *fiber.Ctx, statusCode int, message string, data any, meta any) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends a standardized error JSON response
func Error(c *fiber.Ctx, statusCode int, message string, err any) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// HandleError maps domain errors to standard HTTP status codes and responds
func HandleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	statusCode := GetStatusCode(err)
	if statusCode >= 500 {
		logrus.Error(err)
	}

	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: err.Error(),
	})
}

// GetStatusCode returns the HTTP status code based on the domain error
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrBadParamInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, domain.ErrInternalServerError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
