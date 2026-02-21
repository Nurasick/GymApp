package domain

import "context"

type Recipe struct {
	ID          int64
	UserID      int64
	Ingredients string
	AIResponse  string
	CreatedAt   int64
}

type RecipeRepository interface {
	Create(ctx context.Context, recipe *Recipe) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*Recipe, error)
	GetByID(ctx context.Context, id, userID int64) (*Recipe, error)
}

type RecipeService interface {
	GenerateFromText(ctx context.Context, userID int64, ingredients string) (*Recipe, error)
	GenerateFromImage(ctx context.Context, userID int64, imageData []byte) (*Recipe, error)
	GetHistory(ctx context.Context, userID int64, limit, offset int) ([]*Recipe, error)
}
