package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksysoev/tg-feeder/pkg/core"
)

// newMessage constructs a Telegram message configuration with optional inline keyboard buttons based on given responses.
func newMessage(chatID int64, r *core.Response) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, r.Message)

	return msg
}

// newTextMessage constructs a Telegram message configuration with text and removes the keyboard from the chat.
func newTextMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	return msg
}
