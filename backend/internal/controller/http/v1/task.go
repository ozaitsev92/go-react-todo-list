package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

type taskRoutes struct {
	tasks *usecase.TaskUseCase
	log   *slog.Logger
}

type taskTextRequest struct {
	Text string `json:"text" binding:"required"`
}

func newTaskRoutes(g *gin.RouterGroup, tasks *usecase.TaskUseCase, log *slog.Logger) {
	r := &taskRoutes{tasks: tasks, log: log}

	g.GET("", r.list)
	g.POST("", r.create)
	g.GET("/:id", r.get)
	g.PUT("/:id", r.update)
	g.DELETE("/:id", r.delete)
	g.PUT("/:id/mark-completed", r.setCompleted(true))
	g.PUT("/:id/mark-not-completed", r.setCompleted(false))
}

func (r *taskRoutes) list(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	tasks, err := r.tasks.GetAllTasksForUser(c.Request.Context(), u.ID)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTaskCollection(tasks))
}

func (r *taskRoutes) create(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	var req taskTextRequest
	if !bindJSON(c, &req) {
		return
	}

	t, err := r.tasks.CreateTask(c.Request.Context(), req.Text, u.ID)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	c.Header("Location", "/v1/tasks/"+t.ID.String())
	c.JSON(http.StatusCreated, model.ToResponseFromTask(t))
}

func (r *taskRoutes) get(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	id, ok := pathID(c)
	if !ok {
		return
	}

	t, err := r.tasks.GetTask(c.Request.Context(), id, u.ID)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(t))
}

func (r *taskRoutes) update(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	id, ok := pathID(c)
	if !ok {
		return
	}

	var req taskTextRequest
	if !bindJSON(c, &req) {
		return
	}

	t, err := r.tasks.UpdateTask(c.Request.Context(), id, req.Text, u.ID)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromTask(t))
}

func (r *taskRoutes) delete(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	id, ok := pathID(c)
	if !ok {
		return
	}

	if err := r.tasks.DeleteTask(c.Request.Context(), id, u.ID); err != nil {
		writeError(c, r.log, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *taskRoutes) setCompleted(completed bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUser(c)
		if !ok {
			return
		}

		id, ok := pathID(c)
		if !ok {
			return
		}

		t, err := r.tasks.SetTaskCompleted(c.Request.Context(), id, completed, u.ID)
		if err != nil {
			writeError(c, r.log, err)
			return
		}

		c.JSON(http.StatusOK, model.ToResponseFromTask(t))
	}
}
