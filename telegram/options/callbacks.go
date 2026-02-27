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

type callbackQueryParams struct {
	prefix string
	call   handlers.HandlerFunc
}

func (params *callbackQueryParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithCallbackQueryDataHandler(
		params.prefix, bot.MatchTypePrefix, adapter(params.call),
	)
}

type callbackQueryParamsBuilder struct {
	prefix *string
	call   handlers.HandlerFunc
}

func newCallbackQueryParamsBuilder() *callbackQueryParamsBuilder {
	return &callbackQueryParamsBuilder{}
}

func (builder *callbackQueryParamsBuilder) withPrefix(prefix string) *callbackQueryParamsBuilder {
	builder.prefix = &prefix
	return builder
}

func (builder *callbackQueryParamsBuilder) withCall(call handlers.HandlerFunc) *callbackQueryParamsBuilder {
	builder.call = call
	return builder
}

func (builder *callbackQueryParamsBuilder) build() OptionParams {
	return &callbackQueryParams{
		prefix: *builder.prefix,
		call:   builder.call,
	}
}

func NewNotificationInfoCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("n_").
		withCall(notificationInfo).
		build()
}

func notificationInfo(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "n_", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.Services.Notification.GetByID(s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf(
		"<b>%v</b> — %v <b>$%v</b>",
		n.Symbol, n.Sign.When(), utils.FloatComma(n.Amount),
	)
	kb := view.BuildNotificationInfoKeyboard(n)
	h.SendMessageHTMLWithMarkup(ctx, text, kb)
}

func NewRequestDeleteNotificationCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("rdn_").
		withCall(requestDeleteNotification).
		build()
}

func requestDeleteNotification(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "rdn_", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.Services.Notification.GetByID(s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf("Delete <b>%v %v $%v</b> alert?", n.Symbol, n.Sign, utils.FloatComma(n.Amount))
	kb := view.BuildConfirmDeleteNotificationKeyboard(n)
	h.SendMessageHTMLWithMarkup(ctx, text, kb)
}

func NewDeleteNotificationCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("dn_").
		withCall(deleteNotification).
		build()
}

func deleteNotification(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "dn_", "", 1)
	s = strings.Trim(s, " ")

	if err := h.Services.Notification.DeleteByID(s); err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed delete notification", err)
		return
	}

	h.SendMessage(ctx, "Alert deleted.")
}

func NewDeleteMessageCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("dm_").
		withCall(deleteMessage).
		build()
}

func deleteMessage(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	h.DeleteMessage(ctx, update.CallbackQuery.Message.Message.ID)
}
