package repository

import (
	"context"
	"short-links/internal/links/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db,
	}
}

func (l *LinkRepo) CreateLink(ctx context.Context, userID int, link, code string) (dto.Link, error) {
	query := `INSERT INTO links (original_link, code, user_id) VALUES ($1, $2, $3) RETURNING id, code, created_at`

	var result dto.Link

	err := l.db.QueryRow(
		ctx,
		query,
		link,
		code,
		userID,
	).Scan(
		&result.ID,
		&result.Code,
		&result.CreatedAt,
	)

	if err != nil {
		return dto.Link{}, err
	}

	return result, nil
}

func (l *LinkRepo) GetLinks(ctx context.Context, userID int) ([]dto.Link, error) {
	query := `SELECT id, code, created_at FROM links WHERE user_id = $1`

	rows, err := l.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dto.Link, 0)

	for rows.Next() {
		var link dto.Link

		if err := rows.Scan(&link.ID, &link.Code, &link.CreatedAt); err != nil {
			return nil, err
		}

		result = append(result, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
