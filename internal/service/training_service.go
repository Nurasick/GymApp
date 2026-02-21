package service

import (
	"context"
	"fmt"

	"gymapp/internal/domain"
)

type TrainingService struct {
	trainingRepo domain.TrainingRepository
	userRepo     domain.UserRepository
	aiService    *AIService
}

func NewTrainingService(
	trainingRepo domain.TrainingRepository,
	userRepo domain.UserRepository,
	aiService *AIService,
) *TrainingService {
	return &TrainingService{
		trainingRepo: trainingRepo,
		userRepo:     userRepo,
		aiService:    aiService,
	}
}

func (s *TrainingService) GeneratePlan(ctx context.Context, userID int64, req *domain.GeneratePlanRequest) (*domain.TrainingPlan, error) {
	if req.AvailableDays <= 0 || req.AvailableDays > 7 {
		return nil, fmt.Errorf("available_days must be between 1 and 7")
	}

	if req.TargetWeight <= 0 {
		return nil, fmt.Errorf("target_weight must be positive")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	planJSON, err := s.aiService.GenerateTrainingPlan(ctx, user.Weight, req.TargetWeight, user.Height, req.AvailableDays)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	plan := &domain.TrainingPlan{
		UserID:   userID,
		PlanJSON: planJSON,
	}

	if err := s.trainingRepo.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to save plan: %w", err)
	}

	return plan, nil
}

func (s *TrainingService) GetLatest(ctx context.Context, userID int64) (*domain.TrainingPlan, error) {
	return s.trainingRepo.GetLatestByUserID(ctx, userID)
}
