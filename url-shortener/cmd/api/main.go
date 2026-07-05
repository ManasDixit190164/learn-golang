package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/manasdixit/url-shortener/internal/config"
	"github.com/manasdixit/url-shortener/internal/database"
	"github.com/manasdixit/url-shortener/internal/handler"
	"github.com/manasdixit/url-shortener/internal/repository"
	"github.com/manasdixit/url-shortener/internal/router"
	"github.com/manasdixit/url-shortener/internal/service"
	"github.com/manasdixit/url-shortener/internal/utils"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	jwtManager := utils.NewJWTManager(cfg.JWTAccessSecret, cfg.AccessTokenExpiry)

	userRepo := repository.NewPostgresUserRepository(db)
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(db)
	urlRepo := repository.NewPostgresURLRepository(db)
	clickRepo := repository.NewPostgresClickRepository(db)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager, cfg.RefreshTokenExpiry)
	urlService := service.NewURLService(urlRepo, clickRepo, cfg.BaseURL, cfg.ShortCodeLength)

	authHandler := handler.NewAuthHandler(authService)
	urlHandler := handler.NewURLHandler(urlService, logger)

	r := router.New(router.Dependencies{
		AuthHandler: authHandler,
		URLHandler:  urlHandler,
		JWTManager:  jwtManager,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server started", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
