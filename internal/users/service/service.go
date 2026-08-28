package service

import (
	"context"
	"errors"
	"short-links/internal/users/dto"
	"short-links/internal/utils"
	"strings"
)

type UsersRepo interface {
	CreateUser(ctx context.Context, username, passwordHash string) (dto.User, error)
	GetUserByUsername(ctx context.Context, username string) (dto.User, error)
}

type UsersService struct {
	repo       UsersRepo
	jwtmanager utils.JWTManagerInterface
}

func NewUsersService(repo UsersRepo, jwtmanager utils.JWTManagerInterface) *UsersService {
	return &UsersService{
		repo,
		jwtmanager,
	}
}

func (u *UsersService) RegisterUser(ctx context.Context, user dto.User) (string, error) {
	usernameTrimmed := strings.TrimSpace(user.Username)
	passwordTrimmed := strings.TrimSpace(user.Password)

	if len(passwordTrimmed) < 8 {
		return "", errors.New("password must be more than 7 symbols")
	}

	if len(usernameTrimmed) < 4 {
		return "", errors.New("username must be more than 4 symbols")
	}

	passwordHash, err := utils.HashPassword(passwordTrimmed)
	if err != nil {
		return "", err
	}

	newUser, err := u.repo.CreateUser(ctx, usernameTrimmed, string(passwordHash))
	if err != nil {
		return "", err
	}

	token, err := u.jwtmanager.GenerateTokens(newUser.ID)

	return token, err
}

func (u *UsersService) LoginUser(ctx context.Context, user dto.User) (string, error) {
	usernameTrimmed := strings.TrimSpace(user.Username)
	passwordTrimmed := strings.TrimSpace(user.Password)

	if len(passwordTrimmed) < 8 {
		return "", errors.New("password must be more than 7 symbols")
	}

	if len(usernameTrimmed) < 4 {
		return "", errors.New("username must be more than 4 symbols")
	}

	existedUser, err := u.repo.GetUserByUsername(ctx, usernameTrimmed)
	if err != nil {
		return "", err
	}

	err = utils.CompareHashAndPassword([]byte(existedUser.Password), []byte(passwordTrimmed))
	if err != nil {
		return "", err
	}

	token, err := u.jwtmanager.GenerateTokens(existedUser.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
