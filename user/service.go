package user

import (
	"context"
	"net/mail"
	"time"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/password"
)

type Service struct {
	userRepo domain.UserRepository
}

// NewService creates a new user service instance
func NewService(u domain.UserRepository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (s *Service) Fetch(ctx context.Context, cursor string, num int64) ([]domain.User, string, error) {
	if num <= 0 {
		num = 10
	}
	return s.userRepo.Fetch(ctx, cursor, num)
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *Service) Store(ctx context.Context, u *domain.User) error {
	if u.Email == "" || u.Password == "" {
		return domain.ErrBadParamInput
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return domain.ErrBadParamInput
	}

	if len(u.Password) < 6 {
		return domain.ErrBadParamInput
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, u.Email)
	if err == nil && existing.ID != 0 {
		return domain.ErrConflict
	}

	hashed, err := password.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	return s.userRepo.Store(ctx, u)
}

func (s *Service) Update(ctx context.Context, u *domain.User) error {
	existing, err := s.userRepo.GetByID(ctx, u.ID)
	if err != nil {
		return err
	}

	if u.Email != "" && u.Email != existing.Email {
		if _, err := mail.ParseAddress(u.Email); err != nil {
			return domain.ErrBadParamInput
		}
		byEmail, err := s.userRepo.GetByEmail(ctx, u.Email)
		if err == nil && byEmail.ID != u.ID {
			return domain.ErrConflict
		}
		existing.Email = u.Email
	}

	if u.Password != "" {
		if len(u.Password) < 6 {
			return domain.ErrBadParamInput
		}
		hashed, err := password.Hash(u.Password)
		if err != nil {
			return err
		}
		existing.Password = hashed
	}

	existing.IsActive = u.IsActive
	existing.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, &existing)
	if err != nil {
		return err
	}

	*u = existing
	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *Service) SetActive(ctx context.Context, id int64, isActive bool) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepo.SetActive(ctx, id, isActive)
}

