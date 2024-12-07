package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// todo: refactor. too many params
func NewRouter(handler *gin.Engine, cfg config.Config, l logger.Interface, jwtService *jwt.JWTService, t *usecase.TaskUseCase, u *usecase.UserUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newTaskRoutes(h, l, jwtService, u, t)
		newUserRoutes(h, l, jwtService, u)
	}
}
