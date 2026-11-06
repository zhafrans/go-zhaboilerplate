package auth

import (
	"context"
	"net/mail"
	"time"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/password"
	"go-zhaboilerplate/internal/token"
)

type Service struct {
	userRepo  domain.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewService creates a new auth service instance
func NewService(u domain.UserRepository, jwtSecret string, jwtExpiry time.Duration) *Service {
	return &Service{
		userRepo:  u,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *Service) Register(ctx context.Context, email, plainPassword string) (domain.User, error) {
	if email == "" || plainPassword == "" {
		return domain.User{}, domain.ErrBadParamInput
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return domain.User{}, domain.ErrBadParamInput
	}

	if len(plainPassword) < 6 {
		return domain.User{}, domain.ErrBadParamInput
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing.ID != 0 {
		return domain.User{}, domain.ErrConflict
	}

	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		return domain.User{}, err
	}

	now := time.Now()
	u := domain.User{
		Email:     email,
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.userRepo.Store(ctx, &u)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (s *Service) Login(ctx context.Context, email, plainPassword string) (domain.LoginResponse, error) {
	if email == "" || plainPassword == "" {
		return domain.LoginResponse{}, domain.ErrBadParamInput
	}

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrUnauthorized
	}

	if !u.IsActive {
		return domain.LoginResponse{}, domain.ErrForbidden
	}

	if !password.Check(u.Password, plainPassword) {
		return domain.LoginResponse{}, domain.ErrUnauthorized
	}

	tok, err := token.GenerateToken(u.ID, u.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		Token: tok,
		User:  u,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, userID int64) (domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

