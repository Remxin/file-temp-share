package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"file-temp-share/backend/internal/domain"
)

type ErrorBody struct {
	Error   string            `json:"error"`
	Details []FieldError      `json:"details,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Abort(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorBody{Error: "not found"})
	case errors.Is(err, domain.ErrConflict):
		c.AbortWithStatusJSON(http.StatusConflict, ErrorBody{Error: "conflict"})
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorBody{Error: "unauthorized"})
	case errors.Is(err, domain.ErrForbidden), errors.Is(err, domain.ErrUserInactive):
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorBody{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidResetToken):
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: "invalid or expired reset token"})
	case errors.Is(err, domain.ErrTransferExpired):
		c.AbortWithStatusJSON(http.StatusGone, ErrorBody{Error: "transfer expired"})
	case errors.Is(err, domain.ErrTransferBlocked):
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorBody{Error: "transfer blocked"})
	case errors.Is(err, domain.ErrFileTooLarge):
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, ErrorBody{Error: "file too large"})
	case errors.Is(err, domain.ErrEmptyFile):
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: "empty file"})
	case errors.Is(err, domain.ErrInvalidTransferPassword):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorBody{Error: "invalid transfer password"})
	case errors.Is(err, domain.ErrTransferPasswordRequired):
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorBody{Error: "transfer password required"})
	case errors.Is(err, domain.ErrDownloadLimitReached):
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorBody{Error: "download limit reached"})
	default:
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			details := make([]FieldError, 0, len(ve))
			for _, fe := range ve {
				msg := fe.Tag()
				fields[fe.Field()] = msg
				details = append(details, FieldError{Field: fe.Field(), Message: msg})
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{
				Error:   "validation failed",
				Fields:  fields,
				Details: details,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorBody{Error: "internal server error"})
	}
}
