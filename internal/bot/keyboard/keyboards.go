package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var MainMenuKeyboard = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
	tgbotapi.NewKeyboardButton("Download sound"),
))
