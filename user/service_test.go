package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/user"
	"go-zhaboilerplate/user/mocks"
)

func TestUserFetch(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockUsers := []domain.User{
		{ID: 1, Email: "a@example.com"},
	}

	mockRepo.On("Fetch", mock.Anything, "", int64(10)).Return(mockUsers, "next", nil)

	s := user.NewService(mockRepo)
	res, next, err := s.Fetch(context.TODO(), "", 10)

	assert.NoError(t, err)
	assert.Equal(t, "next", next)
	assert.Len(t, res, 1)
}

func TestUserGetByID(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockUser := domain.User{ID: 1, Email: "a@example.com"}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(mockUser, nil)

	s := user.NewService(mockRepo)
	res, err := s.GetByID(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "a@example.com", res.Email)
}

func TestUserStore(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		newUser := domain.User{
			Email:    "new@example.com",
			Password: "password123",
			IsActive: true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(domain.User{}, domain.ErrNotFound)
		mockRepo.On("Store", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == "new@example.com" && u.Password != "password123" // password hashed
		})).Return(nil)

		s := user.NewService(mockRepo)
		err := s.Store(context.TODO(), &newUser)

		assert.NoError(t, err)
	})

	t.Run("conflict-email", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		newUser := domain.User{
			Email:    "exist@example.com",
			Password: "password123",
		}

		mockRepo.On("GetByEmail", mock.Anything, "exist@example.com").Return(domain.User{ID: 2, Email: "exist@example.com"}, nil)

		s := user.NewService(mockRepo)
		err := s.Store(context.TODO(), &newUser)

		assert.Equal(t, domain.ErrConflict, err)
	})

	t.Run("invalid-email", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		newUser := domain.User{
			Email:    "not-an-email",
			Password: "password123",
		}

		s := user.NewService(mockRepo)
		err := s.Store(context.TODO(), &newUser)

		assert.Equal(t, domain.ErrBadParamInput, err)
	})
}

func TestUserUpdate(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	existing := domain.User{ID: 1, Email: "old@example.com", Password: "oldhash", IsActive: true}
	updated := domain.User{ID: 1, Email: "new@example.com", IsActive: false}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)
	mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(domain.User{}, domain.ErrNotFound)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	s := user.NewService(mockRepo)
	err := s.Update(context.TODO(), &updated)

	assert.NoError(t, err)
	assert.False(t, updated.IsActive)
}

func TestUserDelete(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	existing := domain.User{ID: 1, Email: "old@example.com"}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	s := user.NewService(mockRepo)
	err := s.Delete(context.TODO(), 1)

	assert.NoError(t, err)
}

func TestUserSetActive(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	existing := domain.User{ID: 1, Email: "old@example.com", IsActive: false, UpdatedAt: time.Now()}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)
	mockRepo.On("SetActive", mock.Anything, int64(1), true).Return(nil)

	s := user.NewService(mockRepo)
	err := s.SetActive(context.TODO(), 1, true)

	assert.NoError(t, err)
}

