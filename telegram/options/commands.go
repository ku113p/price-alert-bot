package options

import (
	"context"
	"github.com/ku113p/price-alert-bot/telegram/handlers"
	"github.com/ku113p/price-alert-bot/telegram/helpers"
	"github.com/ku113p/price-alert-bot/telegram/services"
	"github.com/ku113p/price-alert-bot/telegram/view"
	"github.com/ku113p/price-alert-bot/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type commandParams struct {
	name string
	call handlers.HandlerFunc
}

func (params *commandParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithMessageTextHandler(
		params.name, bot.MatchTypeCommand, adapter(params.call),
	)
}

type commandParamsBuilder struct {
	name *string
	call handlers.HandlerFunc
}

func newCommandParamsBuilder() *commandParamsBuilder {
	return &commandParamsBuilder{}
}

func (builder *commandParamsBuilder) withName(name string) *commandParamsBuilder {
	builder.name = &name
	return builder
}

func (builder *commandParamsBuilder) withCall(call handlers.HandlerFunc) *commandParamsBuilder {
	builder.call = call
	return builder
}

func (builder *commandParamsBuilder) build() *commandParams {
	return &commandParams{
		name: *builder.name,
		call: builder.call,
	}
}

func NewStartCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("start").
		withCall(start).
		build()
}

func start(ctx context.Context, _ *models.Update, h *helpers.TelegramRequestHelper) {
	text := "Welcome to <b>Crypto Price Alert Bot</b>!\n\n" +
		"Get notified when a cryptocurrency reaches your target price.\n\n" +
		"Use /add to create your first alert.\n" +
		"Use /help to see all commands."
	h.SendMessageHTML(ctx, text)
}

func NewHelpCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("help").
		withCall(help).
		build()
}

func help(ctx context.Context, _ *models.Update, h *helpers.TelegramRequestHelper) {
	text := "<b>Crypto Price Alert Bot</b>\n\n" +
		"Get notified when a cryptocurrency reaches your target price.\n\n" +
		"<b>Commands:</b>\n" +
		"/add <code>SYMBOL SIGN AMOUNT</code> — create a price alert\n" +
		"/list — view your active alerts\n" +
		"/help — show this message\n\n" +
		"<b>Examples:</b>\n" +
		"<code>/add BTC &gt; 100000</code> — alert when Bitcoin rises above $100,000\n" +
		"<code>/add ETH &lt; 2000</code> — alert when Ethereum drops below $2,000\n\n" +
		"<b>How it works:</b>\n" +
		"1. Create an alert with /add\n" +
		"2. The bot monitors prices in real time\n" +
		"3. You get a notification when the price hits your target"
	h.SendMessageHTML(ctx, text)
}

func NewAddCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("add").
		withCall(add).
		build()
}

func add(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	s := strings.Replace(update.Message.Text, "/add ", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.Services.Notification.Create(h.User, s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}

		h.SendUnexpectedError(ctx, "failed to create notification", err)
		return
	}

	h.SendMessageHTML(ctx, fmt.Sprintf(
		"<b>Alert created</b>\n\n%v %v $%v",
		n.Symbol, n.Sign.HTML(), utils.FloatComma(n.Amount),
	))
}

func NewListCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("list").
		withCall(list).
		build()
}

func list(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	ns, err := h.Services.Notification.GetByUser(h.User)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed get list notifications", err)
		return
	}

	if len(ns) == 0 {
		h.SendMessageHTML(ctx, "You have no active alerts.\n\nUse /add to create one.")
		return
	}

	text := fmt.Sprintf("<b>Your alerts</b> (%d)", len(ns))
	kb := view.BuildNotificationsKeyboard(ns)
	h.SendMessageHTMLWithMarkup(ctx, text, kb)
}
