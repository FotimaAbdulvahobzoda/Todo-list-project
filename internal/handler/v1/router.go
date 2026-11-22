package v1

import (
	"todo_app/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srvs *service.Service
}

func NewHandler(srvs *service.Service) *Handler {
	return &Handler{srvs: srvs}
}

func (h *Handler) InitRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	users := rg.Group("/users", h.authMiddleware)
	{
		users.GET("/me", h.getMe)
		users.GET("", h.listUser)
	}

	tasks := rg.Group("/tasks", h.authMiddleware)
	{
		tasks.POST("", h.createTask)
		tasks.GET("", h.listTask)
		tasks.GET("/:id", h.getTask)
		tasks.PUT("/:id", h.updateTask)
		tasks.PATCH("/:id/completed", h.setTaskCompleted)
		tasks.DELETE("/:id", h.deleteTask)
	}

}
