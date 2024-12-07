package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
)

type taskRoutes struct {
	l          logger.Interface
	jwtService *jwt.JWTService
	u          *usecase.UserUseCase
	t          *usecase.TaskUseCase
}

// todo: refactor. too many params
func newTaskRoutes(handler *gin.RouterGroup, l logger.Interface, jwtService *jwt.JWTService, u *usecase.UserUseCase, t *usecase.TaskUseCase) {
	r := &taskRoutes{l, jwtService, u, t}

	h := handler.Group("/tasks")
	h.Use(middleware.JwtMiddleware(u, jwtService))
	{
		h.GET("", r.index)
		h.POST("", r.createTask)
		h.PUT("/:id", r.updateTask)
		h.DELETE("/:id", r.deleteTask)
		h.PUT("/:id/mark-completed", r.markTaskCompleted)
		h.PUT("/:id/mark-not-completed", r.markTaskNotCompleted)
	}
}

func (r *taskRoutes) index(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	tasks, err := r.t.GetAllTasksForUser(c.Request.Context(), uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - index")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTaskCollection(tasks))
}

func (r *taskRoutes) createTask(c *gin.Context) {
	type createTaskRequest struct {
		Text string `json:"text" binding:"required"`
	}
	var request createTaskRequest

	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - createTask")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})

		return
	}

	task, err := r.t.CreateTask(c.Request.Context(), request.Text, uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - createTask")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(task))
}

func (r *taskRoutes) updateTask(c *gin.Context) {
	type updateTaskRequest struct {
		Text string `json:"text" binding:"required"`
	}
	var request updateTaskRequest

	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - updateTask")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})

		return
	}

	id := c.Param("id")

	task, err := r.t.UpdateTask(c.Request.Context(), uuid.MustParse(id), request.Text, uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - updateTask")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(task))
}

func (r *taskRoutes) deleteTask(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	id := c.Param("id")

	err := r.t.DeleteTask(c.Request.Context(), uuid.MustParse(id), uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - deleteTask")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (r *taskRoutes) markTaskCompleted(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	id := c.Param("id")

	task, err := r.t.MarkTaskCompleted(c.Request.Context(), uuid.MustParse(id), uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - markTaskCompleted")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(task))
}

func (r *taskRoutes) markTaskNotCompleted(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	id := c.Param("id")

	task, err := r.t.MarkTaskNotCompleted(c.Request.Context(), uuid.MustParse(id), uuid.MustParse(userID))
	if err != nil {
		r.l.Error(err, "http - v1 - markTaskNotCompleted")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(task))
}
