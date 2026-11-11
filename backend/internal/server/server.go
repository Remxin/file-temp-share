package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/handler"
	"file-temp-share/backend/internal/middleware"
)

type Dependencies struct {
	AuthHandler  *handler.AuthHandler
	TokenManager *auth.TokenManager
}

// New builds the Gin engine with health, swagger and API routes.
func New(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", deps.AuthHandler.Register)
			authGroup.POST("/login", deps.AuthHandler.Login)
			authGroup.POST("/password-reset", deps.AuthHandler.RequestPasswordReset)
			authGroup.POST("/password-reset/confirm", deps.AuthHandler.ConfirmPasswordReset)
			authGroup.GET("/me", middleware.AuthRequired(deps.TokenManager), deps.AuthHandler.Me)
		}
	}

	return r
}

// health godoc
//
//	@Summary		Health check
//	@Description	Returns service liveness
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
