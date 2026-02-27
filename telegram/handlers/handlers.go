package handlers

import (
	"context"
	"github.com/ku113p/price-alert-bot/app"
	"github.com/ku113p/price-alert-bot/telegram/helpers"
	"github.com/ku113p/price-alert-bot/telegram/middleware"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerFunc func(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper)

type HandlerAdatper func(HandlerFunc) bot.HandlerFunc

func GetAdapter(app *app.App) HandlerAdatper {
	return func(fn HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			user := middleware.ContextUser(ctx)
			telegramHelper := helpers.NewTelegramRequestHelper(bot, user, app)

			fn(ctx, update, telegramHelper)
		}
	}
}
