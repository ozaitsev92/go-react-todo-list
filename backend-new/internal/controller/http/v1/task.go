package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
)

type taskRoutes struct {
	l logger.Interface
	t *usecase.TaskUseCase
}

// todo: too many params
func newTaskRoutes(handler *gin.RouterGroup, l logger.Interface, _ *jwt.JWTService, t *usecase.TaskUseCase) {
	r := &taskRoutes{l, t}

	h := handler.Group("/tasks")
	{
		h.GET("/", r.index)
		h.POST("/", r.createTask)
		h.PUT("/:id", r.updateTask)
		h.DELETE("/:id", r.deleteTask)
		h.PUT("/:id/toggle-completion", r.toggleStatus)
	}
}

func (r *taskRoutes) index(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}

func (r *taskRoutes) createTask(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}

func (r *taskRoutes) updateTask(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}

func (r *taskRoutes) deleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}

func (r *taskRoutes) toggleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}
