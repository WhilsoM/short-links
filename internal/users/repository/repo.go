package repository

import (
	"context"
	"short-links/internal/users/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db,
	}
}

func (u *UsersRepo) CreateUser(ctx context.Context, username, passwordHash string) (dto.User, error) {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username`

	var newUser dto.User

	row := u.db.QueryRow(ctx, query, username, passwordHash)
	if err := row.Scan(&newUser.ID, &newUser.Username); err != nil {
		return dto.User{}, err
	}

	return newUser, nil
}

func (u *UsersRepo) GetUserByUsername(ctx context.Context, username string) (dto.User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`

	var user dto.User

	row := u.db.QueryRow(ctx, query, username)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return dto.User{}, err
	}

	return user, nil

}
