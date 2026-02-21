package domain

import "context"

type TrainingPlan struct {
	ID        int64
	UserID    int64
	PlanJSON  string
	CreatedAt int64
}

type TrainingRepository interface {
	Create(ctx context.Context, plan *TrainingPlan) error
	GetLatestByUserID(ctx context.Context, userID int64) (*TrainingPlan, error)
	GetByID(ctx context.Context, id, userID int64) (*TrainingPlan, error)
}

type TrainingService interface {
	GeneratePlan(ctx context.Context, userID int64, req *GeneratePlanRequest) (*TrainingPlan, error)
	GetLatest(ctx context.Context, userID int64) (*TrainingPlan, error)
}

type GeneratePlanRequest struct {
	TargetWeight  int `json:"target_weight"`
	AvailableDays int `json:"available_days"`
}
