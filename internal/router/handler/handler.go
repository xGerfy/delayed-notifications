package handler

import (
	"context"
	"delayed/internal/entities"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

type Service interface {
	Create(ctx context.Context, message string, sendAt time.Time) (*entities.Notification, error)
	Get(ctx context.Context, id string) (*entities.Notification, error)
	Delete(ctx context.Context, id string) error
}

type notificanionHandler struct {
	service Service
}

func New(service Service) *notificanionHandler {
	return &notificanionHandler{service: service}
}

type CreateNotificationRequest struct {
	Message string    `json:"message" binding:"required"`
	SendAt  time.Time `json:"send_at"`
}

func (h notificanionHandler) CreateNotification(ctx *ginext.Context) {
	var req CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	n, err := h.service.Create(ctx.Request.Context(), req.Message, req.SendAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, n)
}

func (h notificanionHandler) GetNotificationStatus(ctx *ginext.Context) {
	id := ctx.Param("id")

	n, err := h.service.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, n)
}

func (h notificanionHandler) DeleteNotification(ctx *ginext.Context) {
	id := ctx.Param("id")

	if err := h.service.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
