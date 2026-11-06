package domain

import (
	"context"
)

// LoginResponse represents the login payload response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// AuthService represents the authentication usecase contract
type AuthService interface {
	Register(ctx context.Context, email, password string) (User, error)
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	GetProfile(ctx context.Context, userID int64) (User, error)
}

