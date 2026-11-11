package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/dto"
	"file-temp-share/backend/internal/middleware"
	"file-temp-share/backend/internal/platform/httpx"
	"file-temp-share/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register godoc
//
//	@Summary		Register user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.RegisterRequest	true	"credentials"
//	@Success		201		{object}	dto.AuthResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		409		{object}	httpx.ErrorBody
//	@Router			/api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	res, err := h.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.AuthResponse{
		Token:     res.Token,
		ExpiresAt: res.ExpiresAt,
		User:      dto.NewUserResponse(res.User),
	})
}

// Login godoc
//
//	@Summary		Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest	true	"credentials"
//	@Success		200		{object}	dto.AuthResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Failure		403		{object}	httpx.ErrorBody
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	res, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{
		Token:     res.Token,
		ExpiresAt: res.ExpiresAt,
		User:      dto.NewUserResponse(res.User),
	})
}

// Me godoc
//
//	@Summary		Current user
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	dto.UserResponse
//	@Failure		401	{object}	httpx.ErrorBody
//	@Router			/api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	id, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	user, err := h.auth.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// RequestPasswordReset godoc
//
//	@Summary		Request password reset
//	@Description	Always returns success. Email sending is not implemented yet; reset token is logged server-side for development.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.PasswordResetRequest	true	"email"
//	@Success		202		{object}	dto.MessageResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Router			/api/v1/auth/password-reset [post]
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	if err := h.auth.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusAccepted, dto.MessageResponse{
		Message: "if the account exists, password reset instructions were issued",
	})
}

// ConfirmPasswordReset godoc
//
//	@Summary		Confirm password reset
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.PasswordResetConfirmRequest	true	"token and new password"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorBody
//	@Router			/api/v1/auth/password-reset/confirm [post]
func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req dto.PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	if err := h.auth.ConfirmPasswordReset(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		httpx.Abort(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
