package service

import (
	"context"
	"fmt"
	"time"

	"gymapp/internal/config"
	"gymapp/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	cfg              *config.JWTConfig
}

func NewAuthService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	cfg *config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateTokens(userID int64) (string, string, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(s.cfg.AccessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(s.cfg.RefreshTokenTTL).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		UserID: userID,
		Token:  refreshTokenStr,
	}
	if err := s.refreshTokenRepo.Create(context.Background(), rt); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (s *AuthService) ValidateToken(token string) (int64, error) {
	claims := &jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	sub, ok := (*claims)["sub"]
	if !ok {
		return 0, fmt.Errorf("invalid token: missing sub")
	}

	userID, ok := sub.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token: sub is not a number")
	}

	return int64(userID), nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	claims := &jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	sub, ok := (*claims)["sub"]
	if !ok {
		return "", fmt.Errorf("invalid token: missing sub")
	}

	userID, ok := sub.(float64)
	if !ok {
		return "", fmt.Errorf("invalid token: sub is not a number")
	}

	rt, err := s.refreshTokenRepo.GetByToken(context.Background(), refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token not found: %w", err)
	}

	if rt.ExpiresAt < time.Now().Unix() {
		return "", fmt.Errorf("refresh token expired")
	}

	newClaims := jwt.MapClaims{
		"sub": int64(userID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(s.cfg.AccessTokenTTL).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newAccessToken, err := newToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign new token: %w", err)
	}

	return newAccessToken, nil
}
