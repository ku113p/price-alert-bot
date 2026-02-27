package options

import (
	"github.com/ku113p/price-alert-bot/telegram/handlers"

	"github.com/go-telegram/bot"
)

type OptionParams interface {
	ToOption(handlers.HandlerAdatper) bot.Option
}

type OptionParamsBuilder func() OptionParams
