package services

import "github.com/ku113p/price-alert-bot/app"

type Services struct {
	Notification *NotificationService
}

func NewServices(app *app.App) *Services {
	n := newNotificationService(app)

	return &Services{n}
}
