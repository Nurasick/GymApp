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

type TrainingRepository struct {
	pool *pgxpool.Pool
}

func NewTrainingRepository(pool *pgxpool.Pool) *TrainingRepository {
	return &TrainingRepository{pool: pool}
}

func (r *TrainingRepository) Create(ctx context.Context, plan *domain.TrainingPlan) error {
	plan.CreatedAt = time.Now().Unix()

	query := `
		INSERT INTO training_plans (user_id, plan_json, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query, plan.UserID, plan.PlanJSON, plan.CreatedAt).
		Scan(&plan.ID)

	if err != nil {
		return fmt.Errorf("failed to create training plan: %w", err)
	}

	return nil
}

func (r *TrainingRepository) GetLatestByUserID(ctx context.Context, userID int64) (*domain.TrainingPlan, error) {
	query := `
		SELECT id, user_id, plan_json, created_at
		FROM training_plans WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	plan := &domain.TrainingPlan{}
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&plan.ID, &plan.UserID, &plan.PlanJSON, &plan.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no training plan found")
		}
		return nil, fmt.Errorf("failed to get training plan: %w", err)
	}

	return plan, nil
}

func (r *TrainingRepository) GetByID(ctx context.Context, id, userID int64) (*domain.TrainingPlan, error) {
	query := `
		SELECT id, user_id, plan_json, created_at
		FROM training_plans WHERE id = $1 AND user_id = $2
	`

	plan := &domain.TrainingPlan{}
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&plan.ID, &plan.UserID, &plan.PlanJSON, &plan.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("training plan not found")
		}
		return nil, fmt.Errorf("failed to get training plan: %w", err)
	}

	return plan, nil
}
