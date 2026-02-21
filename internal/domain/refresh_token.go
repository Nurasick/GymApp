package domain

import "context"

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt int64
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpiredTokens(ctx context.Context) error
}
