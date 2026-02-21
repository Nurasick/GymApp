package service

import (
	"context"
	"fmt"

	"gymapp/internal/domain"
)

type RecipeService struct {
	recipeRepo domain.RecipeRepository
	aiService  *AIService
}

func NewRecipeService(recipeRepo domain.RecipeRepository, aiService *AIService) *RecipeService {
	return &RecipeService{
		recipeRepo: recipeRepo,
		aiService:  aiService,
	}
}

func (s *RecipeService) GenerateFromText(ctx context.Context, userID int64, ingredients string) (*domain.Recipe, error) {
	if ingredients == "" {
		return nil, fmt.Errorf("ingredients cannot be empty")
	}

	aiResponse, err := s.aiService.GenerateRecipes(ctx, ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipes: %w", err)
	}

	recipe := &domain.Recipe{
		UserID:      userID,
		Ingredients: ingredients,
		AIResponse:  aiResponse,
	}

	if err := s.recipeRepo.Create(ctx, recipe); err != nil {
		return nil, fmt.Errorf("failed to save recipe: %w", err)
	}

	return recipe, nil
}

func (s *RecipeService) GenerateFromImage(ctx context.Context, userID int64, imageData []byte) (*domain.Recipe, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("image data cannot be empty")
	}

	// TODO: Implement image recognition via AI API
	// For now, use a placeholder
	ingredientsList := "Recipe generated from image"

	aiResponse, err := s.aiService.GenerateRecipes(ctx, ingredientsList)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipes: %w", err)
	}

	recipe := &domain.Recipe{
		UserID:      userID,
		Ingredients: ingredientsList,
		AIResponse:  aiResponse,
	}

	if err := s.recipeRepo.Create(ctx, recipe); err != nil {
		return nil, fmt.Errorf("failed to save recipe: %w", err)
	}

	return recipe, nil
}

func (s *RecipeService) GetHistory(ctx context.Context, userID int64, limit, offset int) ([]*domain.Recipe, error) {
	return s.recipeRepo.GetByUserID(ctx, userID, limit, offset)
}
