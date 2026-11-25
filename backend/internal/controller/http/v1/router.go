// Package v1 is the HTTP transport for the v1 API.
package v1

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Deps are the collaborators the router needs.
type Deps struct {
	Log          *slog.Logger
	Tokens       *jwt.Service
	Tasks        *usecase.TaskUseCase
	Users        *usecase.UserUseCase
	MaxBodyBytes int64
	// Ready checks dependencies such as the database; nil means always ready.
	Ready func(ctx context.Context) error
}

// NewRouter builds the HTTP handler for the whole service.
func NewRouter(d Deps) *gin.Engine {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metrics := middleware.NewMetrics(registry)

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(
		middleware.RequestLogger(d.Log),
		metrics.Handler(),
		middleware.Recovery(d.Log),
		middleware.BodyLimit(d.MaxBodyBytes),
	)
	r.NoRoute(func(c *gin.Context) { abortError(c, http.StatusNotFound, "not found") })
	r.NoMethod(func(c *gin.Context) { abortError(c, http.StatusMethodNotAllowed, "method not allowed") })

	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/readyz", func(c *gin.Context) {
		if d.Ready != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			if err := d.Ready(ctx); err != nil {
				d.Log.WarnContext(ctx, "readiness check failed", "err", err)
				abortError(c, http.StatusServiceUnavailable, "not ready")
				return
			}
		}
		c.Status(http.StatusOK)
	})
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

	auth := middleware.Authenticate(d.Tokens, d.Users, d.Log)
	api := r.Group("/v1")
	newTaskRoutes(api.Group("/tasks", auth), d.Tasks, d.Log)
	newUserRoutes(api.Group("/users"), d.Users, d.Tokens, auth, d.Log)

	return r
}
