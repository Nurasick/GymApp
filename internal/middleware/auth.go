package middleware

import (
	"fmt"
	"strings"

	"gymapp/internal/service"

	"github.com/labstack/echo/v4"
)

const (
	authHeaderKey = "Authorization"
	bearerScheme  = "Bearer"
	userIDCtxKey  = "userID"
)

func JWTAuth(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(authHeaderKey)
			if authHeader == "" {
				return echo.NewHTTPError(401, "missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != bearerScheme {
				return echo.NewHTTPError(401, "invalid authorization header format")
			}

			token := parts[1]
			userID, err := authService.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(401, fmt.Sprintf("invalid token: %v", err))
			}

			c.Set(userIDCtxKey, userID)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (int64, error) {
	userID, ok := c.Get(userIDCtxKey).(int64)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
