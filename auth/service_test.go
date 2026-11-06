package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/auth"
	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/password"
	"go-zhaboilerplate/user/mocks"
)

func TestAuthRegister(t *testing.T) {
	secret := "secret-key"
	expiry := time.Hour

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(domain.User{}, domain.ErrNotFound)
		mockRepo.On("Store", mock.Anything, mock.Anything).Return(nil)

		s := auth.NewService(mockRepo, secret, expiry)
		u, err := s.Register(context.TODO(), "new@example.com", "password123")

		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", u.Email)
		assert.True(t, u.IsActive)
	})

	t.Run("email exists", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "exist@example.com").Return(domain.User{ID: 1, Email: "exist@example.com"}, nil)

		s := auth.NewService(mockRepo, secret, expiry)
		_, err := s.Register(context.TODO(), "exist@example.com", "password123")

		assert.Equal(t, domain.ErrConflict, err)
	})
}

func TestAuthLogin(t *testing.T) {
	secret := "secret-key"
	expiry := time.Hour
	hashed, _ := password.Hash("password123")

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockUser := domain.User{
			ID:       1,
			Email:    "user@example.com",
			Password: hashed,
			IsActive: true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(mockUser, nil)

		s := auth.NewService(mockRepo, secret, expiry)
		res, err := s.Login(context.TODO(), "user@example.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)
		assert.Equal(t, mockUser.ID, res.User.ID)
	})

	t.Run("inactive account", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockUser := domain.User{
			ID:       1,
			Email:    "user@example.com",
			Password: hashed,
			IsActive: false,
		}

		mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(mockUser, nil)

		s := auth.NewService(mockRepo, secret, expiry)
		_, err := s.Login(context.TODO(), "user@example.com", "password123")

		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockUser := domain.User{
			ID:       1,
			Email:    "user@example.com",
			Password: hashed,
			IsActive: true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(mockUser, nil)

		s := auth.NewService(mockRepo, secret, expiry)
		_, err := s.Login(context.TODO(), "user@example.com", "wrongpass")

		assert.Equal(t, domain.ErrUnauthorized, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(domain.User{}, domain.ErrNotFound)

		s := auth.NewService(mockRepo, secret, expiry)
		_, err := s.Login(context.TODO(), "notfound@example.com", "password123")

		assert.Equal(t, domain.ErrUnauthorized, err)
	})
}

func TestAuthGetProfile(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockUser := domain.User{ID: 1, Email: "user@example.com", IsActive: true}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(mockUser, nil)

	s := auth.NewService(mockRepo, "sec", time.Hour)
	u, err := s.GetProfile(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Email, u.Email)
}

