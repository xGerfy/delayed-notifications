package router

import (
	"github.com/wb-go/wbf/ginext"
)

type NotificationHandler interface {
	CreateNotification(ctx *ginext.Context)
	GetNotificationStatus(ctx *ginext.Context)
}

func New(h NotificationHandler) *ginext.Engine {
	router := ginext.New("release")

	router.POST("/notify", h.CreateNotification)
	router.GET("/notify/:id", h.GetNotificationStatus)

	return router
}
