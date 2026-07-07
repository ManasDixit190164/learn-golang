package main // package declaration for the module

import ( // start import block
	"context"   // import package
	"log/slog"  // import package
	"net/http"  // import package
	"os"        // import package
	"os/signal" // import package
	"syscall"   // import package
	"time"      // import package

	"github.com/manasdixit/url-shortener/internal/config"     // import package
	"github.com/manasdixit/url-shortener/internal/database"   // import package
	"github.com/manasdixit/url-shortener/internal/handler"    // import package
	"github.com/manasdixit/url-shortener/internal/repository" // import package
	"github.com/manasdixit/url-shortener/internal/router"     // import package
	"github.com/manasdixit/url-shortener/internal/service"    // import package
	"github.com/manasdixit/url-shortener/internal/utils"      // import package
) // end import block or block scope

func main() { // declare function
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})) // declare and initialize variable

	cfg, err := config.Load() // declare and initialize variable
	if err != nil {           // if condition
		logger.Error("failed to load config", "error", err) // execute statement
		os.Exit(1)                                          // execute statement
	} // end block

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // listen for OS signals and shutdown
	defer stop()                                                                           // defer function call

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL) // open a Postgres connection pool
	if err != nil {                                           // if condition
		logger.Error("failed to connect database", "error", err) // execute statement
		os.Exit(1)                                               // execute statement
	} // end block
	defer db.Close() // defer function call

	jwtManager := utils.NewJWTManager(cfg.JWTAccessSecret, cfg.AccessTokenExpiry) // create JWT token manager

	userRepo := repository.NewPostgresUserRepository(db)                 // declare and initialize variable
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(db) // declare and initialize variable
	urlRepo := repository.NewPostgresURLRepository(db)                   // declare and initialize variable
	clickRepo := repository.NewPostgresClickRepository(db)               // declare and initialize variable

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager, cfg.RefreshTokenExpiry) // instantiate authentication service
	urlService := service.NewURLService(urlRepo, clickRepo, cfg.BaseURL, cfg.ShortCodeLength)             // instantiate short URL service

	authHandler := handler.NewAuthHandler(authService)      // instantiate auth HTTP handler
	urlHandler := handler.NewURLHandler(urlService, logger) // instantiate URL HTTP handler

	r := router.New(router.Dependencies{ // build HTTP router and routes
		AuthHandler: authHandler, // execute statement
		URLHandler:  urlHandler,  // execute statement
		JWTManager:  jwtManager,  // execute statement
	}) // close block

	srv := &http.Server{ // configure HTTP server settings
		Addr:         ":" + cfg.Port,   // execute statement
		Handler:      r,                // execute statement
		ReadTimeout:  10 * time.Second, // execute statement
		WriteTimeout: 10 * time.Second, // execute statement
		IdleTimeout:  60 * time.Second, // execute statement
	} // end block

	go func() { // start goroutine
		logger.Info("server started", "port", cfg.Port)                             // execute statement
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { // start HTTP server
			logger.Error("server error", "error", err) // execute statement
			os.Exit(1)                                 // execute statement
		} // end block
	}() // close block

	<-ctx.Done()                        // execute statement
	logger.Info("shutting down server") // execute statement

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // create a context with timeout
	defer cancel()                                                                   // defer function call

	if err := srv.Shutdown(shutdownCtx); err != nil { // stop server gracefully
		logger.Error("server forced to shutdown", "error", err) // execute statement
		os.Exit(1)                                              // execute statement
	} // end block

	logger.Info("server stopped") // execute statement
} // end block
