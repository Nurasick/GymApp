package http

import (
	"net/http"

	"gymapp/internal/domain"
	"gymapp/internal/middleware"
	"gymapp/internal/service"

	"github.com/labstack/echo/v4"
)

type TrainingHandler struct {
	trainingService *service.TrainingService
}

func NewTrainingHandler(trainingService *service.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingService: trainingService}
}

type GeneratePlanRequest struct {
	TargetWeight  int `json:"target_weight"`
	AvailableDays int `json:"available_days"`
}

type PlanResponse struct {
	ID        int64  `json:"id"`
	PlanJSON  string `json:"plan_json"`
	CreatedAt int64  `json:"created_at"`
}

// GeneratePlan godoc
// @Summary Generate training plan
// @Description Generate a personalized training plan using AI
// @ID training-generate
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body GeneratePlanRequest true "Training plan parameters"
// @Success 201 {object} PlanResponse "Training plan generated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /training/generate [post]
func (h *TrainingHandler) GeneratePlan(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req GeneratePlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	plan, err := h.trainingService.GeneratePlan(c.Request().Context(), userID, &domain.GeneratePlanRequest{
		TargetWeight:  req.TargetWeight,
		AvailableDays: req.AvailableDays,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, PlanResponse{
		ID:        plan.ID,
		PlanJSON:  plan.PlanJSON,
		CreatedAt: plan.CreatedAt,
	})
}

// GetLatest godoc
// @Summary Get latest training plan
// @Description Retrieve the user's latest generated training plan
// @ID training-latest
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} PlanResponse "Latest training plan"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "No training plan found"
// @Router /training/latest [get]
func (h *TrainingHandler) GetLatest(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	plan, err := h.trainingService.GetLatest(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no training plan found")
	}

	return c.JSON(http.StatusOK, PlanResponse{
		ID:        plan.ID,
		PlanJSON:  plan.PlanJSON,
		CreatedAt: plan.CreatedAt,
	})
}

func RegisterTrainingRoutes(e *echo.Echo, auth echo.MiddlewareFunc, trainingService *service.TrainingService) {
	handler := NewTrainingHandler(trainingService)

	g := e.Group("/training", auth)
	g.POST("/generate", handler.GeneratePlan)
	g.GET("/latest", handler.GetLatest)
}
