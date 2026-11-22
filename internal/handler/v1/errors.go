package v1

import (
	"errors"
	"log"
	"net/http"
	"todo_app/internal/handler/response"
	"todo_app/internal/models"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, op string, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidInput),
		errors.Is(err, models.ErrPasswordTooShort),
		errors.Is(err, models.ErrPasswordTooLong),
		errors.Is(err, models.ErrInvalidTitle),
		errors.Is(err, models.ErrInvalidDescription),
		errors.Is(err, models.ErrInvalidPriority),
		errors.Is(err, models.ErrInvalidRole),
		errors.Is(err, models.ErrAssigneeNotFound):
		response.Error(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, models.ErrInvalidCredentials),
		errors.Is(err, models.ErrInvalidToken):
		response.Error(c, http.StatusUnauthorized, err.Error())

	case errors.Is(err, models.ErrForbidden):
		response.Error(c, http.StatusForbidden, err.Error())

	case errors.Is(err, models.ErrUserNotFound),
		errors.Is(err, models.ErrTaskNotFound):
		response.Error(c, http.StatusNotFound, err.Error())

	case errors.Is(err, models.ErrEmailAlreadyExists):
		response.Error(c, http.StatusConflict, err.Error())

	default:
		log.Printf("handler %s: %v", op, err)
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
