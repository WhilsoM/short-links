package service

import (
	"context"
	"short-links/internal/users/dto"
	"testing"
)

type fakeJWTManager struct{}

func (f *fakeJWTManager) GenerateTokens(userID int) (string, error) {
	return "fake-token", nil
}

func (f *fakeJWTManager) ParseToken(tokenString string) (int, error) {
	return 1, nil
}

type fakeRepo struct{}

func (f *fakeRepo) CreateUser(
	ctx context.Context,
	username,
	passwordHash string,
) (dto.User, error) {
	return dto.User{
		ID: 1,
	}, nil
}

func (f *fakeRepo) GetUserByUsername(
	ctx context.Context,
	username string,
) (dto.User, error) {
	return dto.User{}, nil
}

func TestRegisterUser(t *testing.T) {
	repo := &fakeRepo{}
	jwtmanager := &fakeJWTManager{}
	tests := []struct {
		name     string
		username string
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty username and password",
			username: "",
			password: "",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "valid user",
			username: "zxcvvvv",
			password: "zxcvvvvv",
			want:     "fake-token",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUsersService(repo, jwtmanager)

			token, err := svc.RegisterUser(context.Background(), dto.User{
				Username: tt.username,
				Password: tt.password,
			})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if token != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, token)
			}
		})
	}

}
