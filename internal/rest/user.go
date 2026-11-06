package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	validator "gopkg.in/go-playground/validator.v9"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/rest/response"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	IsActive *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	IsActive *bool  `json:"is_active"`
}

type SetActiveRequest struct {
	IsActive *bool `json:"is_active" validate:"required"`
}

// UserHandler represents the HTTP handler for users
type UserHandler struct {
	Service   domain.UserService
	validator *validator.Validate
}

const defaultPageSize = 10

// NewUserHandler will initialize the /users endpoints
func NewUserHandler(app fiber.Router, svc domain.UserService, authMiddleware fiber.Handler) {
	handler := &UserHandler{
		Service:   svc,
		validator: validator.New(),
	}

	users := app.Group("/users")
	if authMiddleware != nil {
		users.Use(authMiddleware)
	}

	users.Get("", handler.Fetch)
	users.Post("", handler.Store)
	users.Get("/:id", handler.GetByID)
	users.Put("/:id", handler.Update)
	users.Delete("/:id", handler.Delete)
	users.Patch("/:id/active", handler.SetActive)
}

func (h *UserHandler) Fetch(c *fiber.Ctx) error {
	numS := c.Query("num")
	num, err := strconv.ParseInt(numS, 10, 64)
	if err != nil || num <= 0 {
		num = defaultPageSize
	}

	cursor := c.Query("cursor")
	ctx := c.UserContext()

	users, nextCursor, err := h.Service.Fetch(ctx, cursor, num)
	if err != nil {
		return response.HandleError(c, err)
	}

	meta := response.Meta{
		Cursor: nextCursor,
		Limit:  num,
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Users fetched successfully", users, meta)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, domain.ErrBadParamInput)
	}

	ctx := c.UserContext()
	u, err := h.Service.GetByID(ctx, id)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "User fetched successfully", u)
}

func (h *UserHandler) Store(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Failed to parse request payload", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err.Error())
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	u := domain.User{
		Email:    req.Email,
		Password: req.Password,
		IsActive: isActive,
	}

	ctx := c.UserContext()
	if err := h.Service.Store(ctx, &u); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "User created successfully", u)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, domain.ErrBadParamInput)
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Failed to parse request payload", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err.Error())
	}

	u := domain.User{
		ID:       id,
		Email:    req.Email,
		Password: req.Password,
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	ctx := c.UserContext()
	if err := h.Service.Update(ctx, &u); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "User updated successfully", u)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, domain.ErrBadParamInput)
	}

	ctx := c.UserContext()
	if err := h.Service.Delete(ctx, id); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) SetActive(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, domain.ErrBadParamInput)
	}

	var req SetActiveRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Failed to parse request payload", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error", err.Error())
	}

	ctx := c.UserContext()
	if err := h.Service.SetActive(ctx, id, *req.IsActive); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "User status updated successfully", map[string]interface{}{
		"id":        id,
		"is_active": *req.IsActive,
	})
}
