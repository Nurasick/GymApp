package http

import (
	"net/http"
	"strconv"

	"gymapp/internal/middleware"
	"gymapp/internal/service"

	"github.com/labstack/echo/v4"
)

type RecipeHandler struct {
	recipeService *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

type RecipeRequest struct {
	Ingredients string `json:"ingredients"`
}

type RecipeResponse struct {
	ID          int64  `json:"id"`
	Ingredients string `json:"ingredients"`
	AIResponse  string `json:"ai_response"`
	CreatedAt   int64  `json:"created_at"`
}

// GenerateFromText godoc
// @Summary Generate recipe from ingredients text
// @Description Generate a recipe using AI based on ingredient list
// @ID recipe-generate-text
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body RecipeRequest true "Ingredients for recipe"
// @Success 201 {object} RecipeResponse "Recipe generated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /recipes/from-text [post]
func (h *RecipeHandler) GenerateFromText(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req RecipeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	recipe, err := h.recipeService.GenerateFromText(c.Request().Context(), userID, req.Ingredients)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, RecipeResponse{
		ID:          recipe.ID,
		Ingredients: recipe.Ingredients,
		AIResponse:  recipe.AIResponse,
		CreatedAt:   recipe.CreatedAt,
	})
}

// GenerateFromImage godoc
// @Summary Generate recipe from image
// @Description Generate a recipe by uploading a food image
// @ID recipe-generate-image
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param image formData file true "Image file"
// @Success 201 {object} RecipeResponse "Recipe generated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /recipes/from-image [post]
func (h *RecipeHandler) GenerateFromImage(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "image file is required")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to open file")
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file")
	}

	recipe, err := h.recipeService.GenerateFromImage(c.Request().Context(), userID, buf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, RecipeResponse{
		ID:          recipe.ID,
		Ingredients: recipe.Ingredients,
		AIResponse:  recipe.AIResponse,
		CreatedAt:   recipe.CreatedAt,
	})
}

// GetHistory godoc
// @Summary Get recipe history
// @Description Retrieve user's recipe generation history
// @ID recipe-history
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Number of records" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} RecipeResponse "Recipe history"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /recipes/history [get]
func (h *RecipeHandler) GetHistory(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	recipes, err := h.recipeService.GetHistory(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []RecipeResponse
	for _, recipe := range recipes {
		response = append(response, RecipeResponse{
			ID:          recipe.ID,
			Ingredients: recipe.Ingredients,
			AIResponse:  recipe.AIResponse,
			CreatedAt:   recipe.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func RegisterRecipeRoutes(e *echo.Echo, auth echo.MiddlewareFunc, recipeService *service.RecipeService) {
	handler := NewRecipeHandler(recipeService)

	g := e.Group("/recipes", auth)
	g.POST("/from-text", handler.GenerateFromText)
	g.POST("/from-image", handler.GenerateFromImage)
	g.GET("/history", handler.GetHistory)
}
