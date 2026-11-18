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
	AdminHandler    *handler.AdminHandler
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
			authGroup.POST("/password", middleware.AuthRequired(deps.TokenManager), deps.AuthHandler.ChangePassword)
		}

		transfers := v1.Group("/transfers")
		transfers.Use(middleware.AuthRequired(deps.TokenManager))
		{
			transfers.GET("", deps.TransferHandler.List)
			transfers.POST("", deps.TransferHandler.Create)
			transfers.GET("/:id", deps.TransferHandler.Get)
			transfers.POST("/:id/files", deps.TransferHandler.UploadFile)
			transfers.GET("/:id/files/:fileId", deps.TransferHandler.DownloadOwnedFile)
		}

		publicTransfers := v1.Group("/public/transfers")
		{
			publicTransfers.GET("/:shortCode", deps.TransferHandler.GetPublic)
			publicTransfers.POST("/:shortCode/unlock", deps.TransferHandler.UnlockPublic)
			publicTransfers.GET("/:shortCode/files/:fileId", deps.TransferHandler.DownloadPublicFile)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(deps.TokenManager), middleware.AdminRequired())
		{
			admin.GET("/users", deps.AdminHandler.ListUsers)
			admin.PATCH("/users/:id", deps.AdminHandler.SetUserActive)
			admin.GET("/transfers", deps.AdminHandler.ListTransfers)
			admin.PATCH("/transfers/:id", deps.AdminHandler.SetTransferStatus)
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
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Transfer-Password")
			c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
