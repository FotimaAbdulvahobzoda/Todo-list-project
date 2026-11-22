package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"todo_app/internal/handler/response"
	"todo_app/internal/models"

	"github.com/gin-gonic/gin"
)

type taskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    *int       `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  *int64     `json:"assignee_id"`
}

func (t taskRequest) toInput() models.CreateTaskInput {
	return models.CreateTaskInput{
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority,
		Deadline:    t.Deadline,
		AssigneeID:  t.AssigneeID,
	}
}

type completeRequest struct {
	Completed *bool `json:"completed" binding:"required"`
}

type TaskResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	OwnerID     int64      `json:"owner_id"`
	AssigneeID  *int64     `json:"assignee_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toTaskResponse(t models.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    t.Priority,
		Deadline:    t.Deadline,
		OwnerID:     t.OwnerID,
		AssigneeID:  t.AssigneeID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func (h *Handler) createTask(c *gin.Context) {
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	task, err := h.srvs.TaskService.CreateTask(c.Request.Context(), currentUserID(c), req.toInput())
	if err != nil {
		writeError(c, "create task", err)
		return
	}

	c.JSON(http.StatusCreated, toTaskResponse(*task))
}

func (h *Handler) listTask(c *gin.Context) {
	filters := models.TaskFilter{
		UserID: currentUserID(c),
		Role:   c.Query("role"),
	}

	slog.Info("list tasks request",
		"user_id", filters.UserID,
		"role", c.Query("role"),
	)

	if value := c.Query("completed"); value != "" {
		complted, err := strconv.ParseBool(value)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "completed must be true or false")
			return
		}
		filters.Completed = &complted
	}

	tasks, err := h.srvs.TaskService.ListTask(c.Request.Context(), filters)
	if err != nil {
		writeError(c, "list tasks", err)
		return
	}

	result := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, toTaskResponse(t))
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) getTask(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}
	task, err := h.srvs.TaskService.GetTask(c.Request.Context(), currentUserID(c), taskID)
	if err != nil {
		writeError(c, "get task", err)
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(*task))
}

func (h *Handler) updateTask(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	task, err := h.srvs.TaskService.UpdateTask(c.Request.Context(), currentUserID(c), taskID, req.toInput())
	if err != nil {
		writeError(c, "update task", err)
		return
	}
	c.JSON(http.StatusOK, toTaskResponse(*task))

}

func (h *Handler) setTaskCompleted(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	var req completeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.srvs.TaskService.SetComplete(c.Request.Context(), currentUserID(c), taskID, *req.Completed)
	if err != nil {
		writeError(c, "set completed", err)
		return
	}
	c.JSON(http.StatusOK, toTaskResponse(*task))
}

func (h *Handler) deleteTask(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.srvs.TaskService.DeleteTask(c.Request.Context(), currentUserID(c), taskID); err != nil {
		writeError(c, "delete task", err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}
