package v1

import (
	"net/http"
	"todo_app/internal/handler/response"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type registerResponse struct {
	ID int64 `json:"id"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := h.srvs.AuthService.Register(c.Request.Context(), req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		writeError(c, "register", err)
		return
	}
	c.JSON(http.StatusCreated, registerResponse{ID: userID})
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
	}
	token, err := h.srvs.AuthService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeError(c, "login", err)
		return
	}
	c.JSON(http.StatusOK, loginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}
