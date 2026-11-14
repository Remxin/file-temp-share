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
	AuthHandler     *handler.AuthHandler
	TransferHandler *handler.TransferHandler
	TokenManager    *auth.TokenManager
}

// New builds the Gin engine with health, swagger and API routes.
func New(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 32 << 20 // 32 MiB; larger parts spill to temp disk
	r.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

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

		transfers := v1.Group("/transfers")
		transfers.Use(middleware.AuthRequired(deps.TokenManager))
		{
			transfers.POST("", deps.TransferHandler.Create)
			transfers.POST("/:id/files", deps.TransferHandler.UploadFile)
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

// corsMiddleware allows the local Next.js app to call the API during development.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch origin {
		case "http://localhost:3000", "http://127.0.0.1:3000":
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
