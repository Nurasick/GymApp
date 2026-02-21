package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gymapp/internal/config"
)

type AIService struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func NewAIService(cfg *config.AIConfig) *AIService {
	return &AIService{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AIService) GenerateRecipes(ctx context.Context, ingredients string) (string, error) {
	if s.apiKey == "" {
		return s.mockRecipeGeneration(ingredients), nil
	}

	prompt := fmt.Sprintf(`You are a helpful nutritionist. Given these ingredients: %s
	
Provide 3 healthy recipe ideas with:
1. Recipe name
2. Ingredients needed (with quantities)
3. Estimated calories
4. Preparation time
5. Health benefits

Format as JSON with array of recipes.`, ingredients)

	return s.callChatAPI(ctx, prompt)
}

func (s *AIService) GenerateTrainingPlan(ctx context.Context, weight, targetWeight, height, availableDays int) (string, error) {
	if s.apiKey == "" {
		return s.mockTrainingPlanGeneration(weight, targetWeight, height, availableDays), nil
	}

	prompt := fmt.Sprintf(`You are an expert fitness coach. Create a personalized training plan with:
Current: Weight=%dkg, Height=%dcm, Available days per week=%d
Goal: Target weight=%dkg

Provide:
1. Weekly workout schedule
2. Exercise type and duration for each day
3. Nutrition guidance
4. Recovery recommendations
5. Expected timeline

Format as JSON.`, weight, height, availableDays, targetWeight)

	return s.callChatAPI(ctx, prompt)
}

func (s *AIService) callChatAPI(ctx context.Context, prompt string) (string, error) {
	req := ChatRequest{
		Model: s.model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		s.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api error: status %d, body: %s", resp.StatusCode, string(data))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (s *AIService) mockRecipeGeneration(ingredients string) string {
	return fmt.Sprintf(`{
  "recipes": [
    {
      "name": "Healthy Salad Bowl with %s",
      "calories": 350,
      "prep_time": "15 minutes",
      "benefits": "Rich in vitamins and fiber"
    }
  ]
}`, ingredients)
}

func (s *AIService) mockTrainingPlanGeneration(weight, targetWeight, height, availableDays int) string {
	return fmt.Sprintf(`{
  "plan": {
    "duration_weeks": 12,
    "weekly_schedule": {
      "monday": "Cardio 30min, Strength 20min",
      "wednesday": "Strength training 45min",
      "friday": "Mixed cardio and agility"
    },
    "nutrition": "Caloric deficit of 500-750 kcal/day",
    "expected_weight_loss": %d
  }
}`, weight-targetWeight)
}
