package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/config"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	taskRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	userRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/httpserver"
	"github.com/ozaitsev92/tododdd/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg config.Config) {
	l := logger.New(cfg)

	// Task Use case
	taskUseCase := usecase.NewTaskUseCase(
		taskRepository.NewRepository(cfg),
	)

	// User Use case
	userUseCase := usecase.NewUserUseCase(
		userRepository.NewRepository(cfg),
	)

	// JWT service
	jwtService := jwt.NewJWTService(
		[]byte(cfg.JWTSigningKey),
		cfg.JWTSessionLength,
		cfg.JWTCookieDomain,
		cfg.JWTSecureCookie,
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, jwtService, taskUseCase, userUseCase)
	httpServer := httpserver.New(
		handler,
		httpserver.Port(cfg.BindAddr),
		httpserver.ReadTimeout(time.Duration(cfg.ReadTimeout)*time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.WriteTimeout)*time.Second),
		httpserver.ShutdownTimeout(time.Duration(cfg.GracefulTimeout)*time.Second),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
