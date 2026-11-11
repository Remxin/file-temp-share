package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"file-temp-share/backend/docs"
	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/config"
	"file-temp-share/backend/internal/handler"
	"file-temp-share/backend/internal/platform/database"
	"file-temp-share/backend/internal/repository"
	"file-temp-share/backend/internal/server"
	"file-temp-share/backend/internal/service"
)

//	@title						file-temp-share API
//	@version					0.1.0
//	@description				Temporary file-sharing API (WeTransfer-like).
//	@BasePath					/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT token. Example: "Bearer {token}"

//	@contact.name	file-temp-share
//	@contact.url	https://github.com/Remxin/file-temp-share

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}
	docs.SwaggerInfo.Host = cfg.SwaggerHost

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	tokenManager := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTTTL)
	userRepo := repository.NewUserRepository(pool)
	resetRepo := repository.NewPasswordResetRepository(pool)
	authService := service.NewAuthService(userRepo, resetRepo, tokenManager)
	authHandler := handler.NewAuthHandler(authService)

	engine := server.New(server.Dependencies{
		AuthHandler:  authHandler,
		TokenManager: tokenManager,
	})

	go func() {
		<-ctx.Done()
		slog.Info("shutting down")
	}()

	slog.Info("api listening",
		"addr", cfg.HTTPAddr,
		"swagger", "http://"+cfg.SwaggerHost+"/swagger/index.html",
		"jwt_ttl", cfg.JWTTTL.String(),
	)
	if err := engine.Run(cfg.HTTPAddr); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
