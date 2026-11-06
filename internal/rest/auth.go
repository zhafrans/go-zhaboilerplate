package rest

import (
	"github.com/gofiber/fiber/v2"
	validator "gopkg.in/go-playground/validator.v9"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/rest/middleware"
	"go-zhaboilerplate/internal/rest/response"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthHandler represents the HTTP handler for authentication
type AuthHandler struct {
	Service   domain.AuthService
	validator *validator.Validate
}

// NewAuthHandler will initialize auth endpoints
func NewAuthHandler(app fiber.Router, svc domain.AuthService, authMiddleware fiber.Handler) {
	handler := &AuthHandler{
		Service:   svc,
		validator: validator.New(),
	}

	auth := app.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	if authMiddleware != nil {
		auth.Get("/me", authMiddleware, handler.Me)
	} else {
		auth.Get("/me", handler.Me)
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Failed to parse request payload", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err.Error())
	}

	ctx := c.UserContext()
	u, err := h.Service.Register(ctx, req.Email, req.Password)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", u)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Failed to parse request payload", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err.Error())
	}

	ctx := c.UserContext()
	res, err := h.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	ctx := c.UserContext()
	u, err := h.Service.GetProfile(ctx, userID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", u)
}
