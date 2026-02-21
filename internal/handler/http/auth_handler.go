package http

import (
	"net/http"

	"gymapp/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserResponse struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Height int    `json:"height,omitempty"`
	Weight int    `json:"weight,omitempty"`
	Goal   string `json:"goal,omitempty"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @ID auth-register
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 201 {object} TokenResponse "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate tokens")
	}

	return c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get tokens
// @ID auth-login
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} TokenResponse "Login successful"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate tokens")
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @ID auth-refresh
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh request"
// @Success 200 {object} map[string]interface{} "New access token"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	accessToken, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   900,
	})
}

func RegisterAuthRoutes(e *echo.Echo, authService *service.AuthService) {
	handler := NewAuthHandler(authService)

	e.POST("/auth/register", handler.Register)
	e.POST("/auth/login", handler.Login)
	e.POST("/auth/refresh", handler.Refresh)
}
