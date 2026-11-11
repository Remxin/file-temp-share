package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/platform/httpx"
)

const (
	ContextUserIDKey = "userID"
	ContextEmailKey  = "userEmail"
	ContextRoleKey   = "userRole"
)

func AuthRequired(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpx.ErrorBody{Error: "unauthorized"})
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := tokens.Parse(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpx.ErrorBody{Error: "unauthorized"})
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}

func UserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ContextUserIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
