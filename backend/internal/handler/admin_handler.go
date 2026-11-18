package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/dto"
	"file-temp-share/backend/internal/middleware"
	"file-temp-share/backend/internal/platform/httpx"
	"file-temp-share/backend/internal/service"
)

type AdminHandler struct {
	admin *service.AdminService
}

func NewAdminHandler(admin *service.AdminService) *AdminHandler {
	return &AdminHandler{admin: admin}
}

// ListUsers godoc
//
//	@Summary		List all users (admin)
//	@Tags			admin
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		dto.UserResponse
//	@Failure		401	{object}	httpx.ErrorBody
//	@Failure		403	{object}	httpx.ErrorBody
//	@Router			/api/v1/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.admin.ListUsers(c.Request.Context())
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.NewUserResponse(u))
	}
	c.JSON(http.StatusOK, out)
}

// SetUserActive godoc
//
//	@Summary		Activate or deactivate a user (admin)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"user id"
//	@Param			body	body		dto.AdminSetUserActiveRequest	true	"active flag"
//	@Success		200		{object}	dto.UserResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Failure		403		{object}	httpx.ErrorBody
//	@Failure		404		{object}	httpx.ErrorBody
//	@Router			/api/v1/admin/users/{id} [patch]
func (h *AdminHandler) SetUserActive(c *gin.Context) {
	actorID, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, domain.ErrValidation)
		return
	}
	var req dto.AdminSetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	user, err := h.admin.SetUserActive(c.Request.Context(), actorID, targetID, *req.IsActive)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// ListTransfers godoc
//
//	@Summary		List all transfers (admin)
//	@Tags			admin
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		dto.AdminTransferListItem
//	@Failure		401	{object}	httpx.ErrorBody
//	@Failure		403	{object}	httpx.ErrorBody
//	@Router			/api/v1/admin/transfers [get]
func (h *AdminHandler) ListTransfers(c *gin.Context) {
	items, err := h.admin.ListTransfers(c.Request.Context())
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	out := make([]dto.AdminTransferListItem, 0, len(items))
	for _, item := range items {
		out = append(out, dto.NewAdminTransferListItem(item, h.admin.ShareURL(item.Transfer)))
	}
	c.JSON(http.StatusOK, out)
}

// SetTransferStatus godoc
//
//	@Summary		Block or unblock a transfer (admin)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string								true	"transfer id"
//	@Param			body	body		dto.AdminSetTransferStatusRequest	true	"status"
//	@Success		200		{object}	dto.TransferResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Failure		403		{object}	httpx.ErrorBody
//	@Failure		404		{object}	httpx.ErrorBody
//	@Router			/api/v1/admin/transfers/{id} [patch]
func (h *AdminHandler) SetTransferStatus(c *gin.Context) {
	transferID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, domain.ErrValidation)
		return
	}
	var req dto.AdminSetTransferStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	transfer, err := h.admin.SetTransferStatus(c.Request.Context(), transferID, req.Status)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewTransferResponse(transfer, h.admin.ShareURL(transfer)))
}
