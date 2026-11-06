package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"go-zhaboilerplate/domain"
)

type AuthService struct {
	mock.Mock
}

func (m *AuthService) Register(ctx context.Context, email, password string) (domain.User, error) {
	ret := m.Called(ctx, email, password)
	var res domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.User); ok {
		res = rf(ctx, email, password)
	} else if ret.Get(0) != nil {
		res = ret.Get(0).(domain.User)
	}
	return res, ret.Error(1)
}

func (m *AuthService) Login(ctx context.Context, email, password string) (domain.LoginResponse, error) {
	ret := m.Called(ctx, email, password)
	var res domain.LoginResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.LoginResponse); ok {
		res = rf(ctx, email, password)
	} else if ret.Get(0) != nil {
		res = ret.Get(0).(domain.LoginResponse)
	}
	return res, ret.Error(1)
}

func (m *AuthService) GetProfile(ctx context.Context, userID int64) (domain.User, error) {
	ret := m.Called(ctx, userID)
	var res domain.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.User); ok {
		res = rf(ctx, userID)
	} else if ret.Get(0) != nil {
		res = ret.Get(0).(domain.User)
	}
	return res, ret.Error(1)
}

