package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gymapp/docs"
	"gymapp/internal/config"
	"gymapp/internal/database"
	httphandler "gymapp/internal/handler/http"
	midauth "gymapp/internal/middleware"
	"gymapp/internal/repository/postgres"
	"gymapp/internal/service"
	"gymapp/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title GymApp API
// @version 1.0
// @description A comprehensive API for managing gym training plans and recipes with AI integration
// @termsOfService http://swagger.io/terms/
// @contact.email support@gymapp.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /
// @schemes http https
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	logger := utils.NewLogger()

	cfg := config.Load()

	fmt.Println("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "")
	fmt.Println("🏋️  GymApp Backend - Startup Sequence")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "" + "\n")

	logger.Infof("📋 Configuration: env=%s, port=%d", cfg.Server.Env, cfg.Server.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Database connection with enhanced logging
	fmt.Println("🔌 Connecting to database...")
	logger.Infof("📍 Database: %s (host=%s, port=%d)", cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port)

	pool, err := database.NewPool(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Errorf("❌ Database connection failed: %v", err)
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify database is accessible
	verifyCtx, verifyCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := pool.Ping(verifyCtx); err != nil {
		verifyCancel()
		logger.Errorf("❌ Failed to ping database: %v", err)
		fmt.Printf("❌ Database ping failed: %v\n", err)
		os.Exit(1)
	}
	verifyCancel()

	fmt.Println("✅ Database connection established")
	logger.Infof("✅ Database connected: %s", cfg.Database.DBName)

	// Check migrations
	var migrationCount int64
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = pool.QueryRow(checkCtx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&migrationCount)
	checkCancel()

	if err == nil && migrationCount > 0 {
		fmt.Printf("✅ Database schema initialized (%d tables found)\n", migrationCount)
		logger.Infof("✅ Database schema ready with %d tables", migrationCount)
	} else {
		fmt.Println("⚠️  No tables found - run 'make migrate-up' to apply migrations")
		logger.Infof("⚠️  No tables found in database - migrations may not be applied")
	}
	fmt.Println()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	recipeRepo := postgres.NewRecipeRepository(pool)
	trainingRepo := postgres.NewTrainingRepository(pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, refreshTokenRepo, &cfg.JWT)
	aiService := service.NewAIService(&cfg.AI)
	recipeService := service.NewRecipeService(recipeRepo, aiService)
	trainingService := service.NewTrainingService(trainingRepo, userRepo, aiService)

	// Setup Echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Register Swagger routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	httphandler.RegisterHealthRoutes(e)
	httphandler.RegisterAuthRoutes(e, authService)

	authMiddleware := midauth.JWTAuth(authService)
	httphandler.RegisterRecipeRoutes(e, authMiddleware, recipeService)
	httphandler.RegisterTrainingRoutes(e, authMiddleware, trainingService)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Infof("shutting down server...")
		fmt.Println("\n⏹️  Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("error during shutdown: %v", err)
			fmt.Printf("Error during shutdown: %v\n", err)
		}

		pool.Close()
	}()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "")
	fmt.Println("🚀 Server Starting")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "")
	fmt.Printf("✅ API Server running at:        http://localhost:%d\n", cfg.Server.Port)
	fmt.Printf("📚 Swagger UI available at:      http://localhost:%d/swagger/index.html\n", cfg.Server.Port)
	fmt.Printf("🏥 Health check endpoint:        http://localhost:%d/health\n", cfg.Server.Port)
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "" + "\n")

	logger.Infof("server started on port %d", cfg.Server.Port)

	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		logger.Errorf("server error: %v", err)
		fmt.Printf("❌ Server error: %v\n", err)
		os.Exit(1)
	}
}
