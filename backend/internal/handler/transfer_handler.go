package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/dto"
	"file-temp-share/backend/internal/middleware"
	"file-temp-share/backend/internal/platform/httpx"
	"file-temp-share/backend/internal/service"
)

type TransferHandler struct {
	transfers      *service.TransferService
	uploadMaxBytes int64
}

func NewTransferHandler(transfers *service.TransferService, uploadMaxBytes int64) *TransferHandler {
	return &TransferHandler{transfers: transfers, uploadMaxBytes: uploadMaxBytes}
}

// Create godoc
//
//	@Summary		Create transfer
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTransferRequest	true	"transfer metadata"
//	@Success		201		{object}	dto.TransferResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Router			/api/v1/transfers [post]
func (h *TransferHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	var req dto.CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, err)
		return
	}
	transfer, err := h.transfers.Create(c.Request.Context(), userID, service.CreateTransferInput{
		Title:          req.Title,
		Message:        req.Message,
		Password:       req.Password,
		ExpiresInHours: req.ExpiresInHours,
		Recipients:     req.Recipients,
	})
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewTransferResponse(transfer, h.transfers.ShareURL(transfer)))
}

// List godoc
//
//	@Summary		List own transfers
//	@Tags			transfers
//	@Produce		json
//	@Security		BearerAuth
//	@Param			status	query		string	false	"active|expired|all"	Enums(active,expired,all)
//	@Success		200		{array}		dto.TransferListItem
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Router			/api/v1/transfers [get]
func (h *TransferHandler) List(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	filter := strings.TrimSpace(c.DefaultQuery("status", "active"))
	items, err := h.transfers.List(c.Request.Context(), userID, filter)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	out := make([]dto.TransferListItem, 0, len(items))
	for _, item := range items {
		base := dto.NewTransferResponse(item.Transfer, h.transfers.ShareURL(item.Transfer))
		out = append(out, dto.TransferListItem{
			TransferResponse: base,
			FileCount:        item.FileCount,
			TotalSizeBytes:   item.TotalSizeBytes,
			RecipientCount:   item.RecipientCount,
			FirstFileName:    item.FirstFileName,
		})
	}
	c.JSON(http.StatusOK, out)
}

// Get godoc
//
//	@Summary		Get own transfer detail
//	@Tags			transfers
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"transfer id"
//	@Success		200	{object}	dto.TransferDetailResponse
//	@Failure		401	{object}	httpx.ErrorBody
//	@Failure		403	{object}	httpx.ErrorBody
//	@Failure		404	{object}	httpx.ErrorBody
//	@Router			/api/v1/transfers/{id} [get]
func (h *TransferHandler) Get(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	transferID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, domain.ErrNotFound)
		return
	}
	detail, err := h.transfers.GetOwned(c.Request.Context(), userID, transferID)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	base := dto.NewTransferResponse(detail.Transfer, h.transfers.ShareURL(detail.Transfer))
	c.JSON(http.StatusOK, dto.TransferDetailResponse{
		TransferResponse: base,
		Files:            dto.NewFileResponses(detail.Files),
		FileCount:        detail.FileCount,
		TotalSizeBytes:   detail.TotalSizeBytes,
		RecipientCount:   detail.RecipientCount,
	})
}

// GetPublic godoc
//
//	@Summary		Public transfer by short code
//	@Tags			public
//	@Produce		json
//	@Param			shortCode	path		string	true	"short code"
//	@Success		200			{object}	dto.PublicTransferResponse
//	@Failure		403			{object}	httpx.ErrorBody
//	@Failure		404			{object}	httpx.ErrorBody
//	@Router			/api/v1/public/transfers/{shortCode} [get]
func (h *TransferHandler) GetPublic(c *gin.Context) {
	shortCode := strings.TrimSpace(c.Param("shortCode"))
	if shortCode == "" {
		httpx.Abort(c, domain.ErrNotFound)
		return
	}
	detail, err := h.transfers.GetPublic(c.Request.Context(), shortCode)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.PublicTransferResponse{
		ShortCode:      detail.Transfer.ShortCode,
		Title:          detail.Transfer.Title,
		Message:        detail.Transfer.Message,
		Status:         detail.Transfer.Status,
		ExpiresAt:      detail.Transfer.ExpiresAt,
		HasPassword:    detail.Transfer.PasswordHash != nil && *detail.Transfer.PasswordHash != "",
		SenderEmail:    detail.SenderEmail,
		Files:          dto.NewFileResponses(detail.Files),
		FileCount:      detail.FileCount,
		TotalSizeBytes: detail.TotalSizeBytes,
	})
}

// UploadFile godoc
//
//	@Summary		Upload file to transfer
//	@Tags			transfers
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"transfer id"
//	@Param			file	formData	file	true	"file"
//	@Success		201		{object}	dto.FileResponse
//	@Failure		400		{object}	httpx.ErrorBody
//	@Failure		401		{object}	httpx.ErrorBody
//	@Failure		403		{object}	httpx.ErrorBody
//	@Failure		404		{object}	httpx.ErrorBody
//	@Failure		410		{object}	httpx.ErrorBody
//	@Failure		413		{object}	httpx.ErrorBody
//	@Router			/api/v1/transfers/{id}/files [post]
func (h *TransferHandler) UploadFile(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		httpx.Abort(c, domain.ErrUnauthorized)
		return
	}
	transferID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, domain.ErrNotFound)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.uploadMaxBytes+1024*1024)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpx.Abort(c, domain.ErrFileTooLarge)
			return
		}
		httpx.Abort(c, domain.ErrValidation)
		return
	}
	if fileHeader.Size <= 0 {
		httpx.Abort(c, domain.ErrEmptyFile)
		return
	}
	if fileHeader.Size > h.uploadMaxBytes {
		httpx.Abort(c, domain.ErrFileTooLarge)
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	defer src.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	file, err := h.transfers.UploadFile(c.Request.Context(), service.UploadFileInput{
		TransferID:   transferID,
		UserID:       userID,
		OriginalName: fileHeader.Filename,
		ContentType:  contentType,
		Size:         fileHeader.Size,
		Body:         src,
	})
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewFileResponse(file))
}
