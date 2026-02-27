package view

import (
	"github.com/ku113p/price-alert-bot/models"
	"github.com/ku113p/price-alert-bot/utils"
	"fmt"

	telegramModels "github.com/go-telegram/bot/models"
)

func BuildNotificationsKeyboard(ns []*models.Notification) *telegramModels.InlineKeyboardMarkup {
	buttons := make([][]telegramModels.InlineKeyboardButton, 0)
	for _, n := range ns {
		row := []telegramModels.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%v %s $%v", n.Symbol, n.Sign, utils.FloatComma(n.Amount)),
				CallbackData: fmt.Sprintf("n_%v", n.ID.String()),
			},
		}
		buttons = append(buttons, row)
	}

	return &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func BuildNotificationInfoKeyboard(n *models.Notification) *telegramModels.InlineKeyboardMarkup {
	return &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegramModels.InlineKeyboardButton{
			{
				{
					Text:         "Delete ❌",
					CallbackData: fmt.Sprintf("rdn_%v", n.ID.String()),
				},
			},
		},
	}
}

func BuildConfirmDeleteNotificationKeyboard(n *models.Notification) *telegramModels.InlineKeyboardMarkup {
	return &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegramModels.InlineKeyboardButton{
			{
				{
					Text:         "Delete ❌",
					CallbackData: fmt.Sprintf("dn_%v", n.ID.String()),
				},
				{
					Text:         "Cancel ⭕",
					CallbackData: fmt.Sprintf("dm_%v", nil),
				},
			},
		},
	}
}
