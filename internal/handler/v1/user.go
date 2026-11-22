package v1

import (
	"net/http"
	"time"
	"todo_app/internal/models"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(u models.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}

func (h *Handler) getMe(c *gin.Context) {
	user, err := h.srvs.UserService.GetByID(c.Request.Context(), currentUserID(c))
	if err != nil {
		writeError(c, "get me", err)
	}

	c.JSON(http.StatusOK, toUserResponse(*user))
}

func (h *Handler) listUser(c *gin.Context) {
	users, err := h.srvs.UserService.GetAll(c.Request.Context())
	if err != nil {
		writeError(c, "list users", err)
	}

	result := make([]userResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toUserResponse(u))
	}
	c.JSON(http.StatusOK, result)
}
