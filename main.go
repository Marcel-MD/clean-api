package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Marcel-MD/clean-api/api"
	"github.com/Marcel-MD/clean-api/api/controllers"
	"github.com/Marcel-MD/clean-api/config"
	"github.com/Marcel-MD/clean-api/data"
	"github.com/Marcel-MD/clean-api/data/repositories"
	"github.com/Marcel-MD/clean-api/services"
	"github.com/rs/zerolog/log"
)

// @title Clean API
// @description This is a sample server for a clean API.
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// User
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, cfg)
	userController := controllers.NewUserController(userService)

	srv := api.NewServer(cfg, userController)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
		log.Info().Msg("All server connections are closed")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Warn().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	if err := data.CloseDB(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to close db connection")
	}

	log.Info().Msg("Server exited properly")
}
