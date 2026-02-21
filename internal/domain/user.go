package domain

import "context"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Height       int
	Weight       int
	Goal         string
	CreatedAt    int64
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GenerateTokens(userID int64) (accessToken, refreshToken string, err error)
	ValidateToken(token string) (userID int64, err error)
	RefreshAccessToken(refreshToken string) (newAccessToken string, err error)
}
