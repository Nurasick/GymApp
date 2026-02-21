package postgres

import (
	"context"
	"fmt"
	"time"

	"gymapp/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecipeRepository struct {
	pool *pgxpool.Pool
}

func NewRecipeRepository(pool *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{pool: pool}
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	recipe.CreatedAt = time.Now().Unix()

	query := `
		INSERT INTO recipes (user_id, ingredients, ai_response, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		recipe.UserID, recipe.Ingredients, recipe.AIResponse, recipe.CreatedAt).
		Scan(&recipe.ID)

	if err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}

	return nil
}

func (r *RecipeRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.Recipe, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
		SELECT id, user_id, ingredients, ai_response, created_at
		FROM recipes WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*domain.Recipe
	for rows.Next() {
		recipe := &domain.Recipe{}
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Ingredients, &recipe.AIResponse, &recipe.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *RecipeRepository) GetByID(ctx context.Context, id, userID int64) (*domain.Recipe, error) {
	query := `
		SELECT id, user_id, ingredients, ai_response, created_at
		FROM recipes WHERE id = $1 AND user_id = $2
	`

	recipe := &domain.Recipe{}
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Ingredients, &recipe.AIResponse, &recipe.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("recipe not found")
		}
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	return recipe, nil
}
