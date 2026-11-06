package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/domain"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) Fetch(ctx context.Context, cursor string, num int64) ([]domain.User, string, error) {
	ret := m.Called(ctx, cursor, num)
	var res []domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) []domain.User); ok {
		res = rf(ctx, cursor, num)
	} else if ret.Get(0) != nil {
		res = ret.Get(0).([]domain.User)
	}
	return res, ret.String(1), ret.Error(2)
}

func (m *UserService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	ret := m.Called(ctx, id)
	var res domain.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.User); ok {
		res = rf(ctx, id)
	} else if ret.Get(0) != nil {
		res = ret.Get(0).(domain.User)
	}
	return res, ret.Error(1)
}

func (m *UserService) Store(ctx context.Context, u *domain.User) error {
	ret := m.Called(ctx, u)
	return ret.Error(0)
}

func (m *UserService) Update(ctx context.Context, u *domain.User) error {
	ret := m.Called(ctx, u)
	return ret.Error(0)
}

func (m *UserService) Delete(ctx context.Context, id int64) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *UserService) SetActive(ctx context.Context, id int64, isActive bool) error {
	ret := m.Called(ctx, id, isActive)
	return ret.Error(0)
}

