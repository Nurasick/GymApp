package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gymapp/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	rt.ExpiresAt = time.Now().AddDate(0, 0, 7).Unix()

	err := r.pool.QueryRow(ctx, query, rt.UserID, rt.Token, rt.ExpiresAt, time.Now().Unix()).
		Scan(&rt.ID)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at
		FROM refresh_tokens WHERE token = $1
	`

	rt := &domain.RefreshToken{}
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return rt, nil
}

func (r *RefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`

	result, err := r.pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`

	_, err := r.pool.Exec(ctx, query, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}
