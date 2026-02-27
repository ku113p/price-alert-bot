package options

import (
	"context"
	"github.com/ku113p/price-alert-bot/telegram/handlers"
	"github.com/ku113p/price-alert-bot/telegram/helpers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type defaultParams struct{}

func (p *defaultParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithDefaultHandler(adapter(defaultEcho))
}

func NewDefaultParams() OptionParams {
	return &defaultParams{}
}

func defaultEcho(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	if update.Message != nil {
		h.SendMessage(ctx, "Unknown command. Use /help to see available commands.")
	}
}
