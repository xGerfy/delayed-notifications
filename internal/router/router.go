package router

import (
	"github.com/wb-go/wbf/ginext"
)

type NotificationHandler interface {
	CreateNotification(ctx *ginext.Context)
	GetNotificationStatus(ctx *ginext.Context)
	DeleteNotification(ctx *ginext.Context)
}

func New(h NotificationHandler) *ginext.Engine {
	router := ginext.New("release")

	router.POST("/notify", h.CreateNotification)
	router.GET("/notify/:id", h.GetNotificationStatus)
	router.DELETE("/notify/:id", h.DeleteNotification)

	return router
}
