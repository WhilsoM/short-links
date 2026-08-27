package service

import (
	"context"
	"short-links/internal/users/dto"
	"testing"
)

type fakeRepo struct {
}

func (f *fakeRepo) CreateUser(ctx context.Context, username, passwordHash string) (dto.User, error) {
	return dto.User{}, nil
}

func (f *fakeRepo) GetUserByUsername(ctx context.Context, username string) (dto.User, error) {
	return dto.User{}, nil
}

func TestRegisterUser(t *testing.T) {
	NewUsersService(fakeRepo, nil)
}
