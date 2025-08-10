package middleware

import (
	"context"
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// WithErrorHandling adds error handling middleware to a Handler.
// It intercepts errors returned by the next Handler and generates an appropriate error message response for the user.
// It uses the localized message printer from the context to create user-friendly error messages.
// Returns a Middleware wrapping the original Handler with error handling logic.
func WithErrorHandling() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
			if message == nil {
				return tgbotapi.MessageConfig{}, errors.New("message is nil")
			}

			msgConfig, err := next.Handle(ctx, message)
			if err != nil {
				var chatID int64
				if message.Chat != nil {
					chatID = message.Chat.ID
				}

				slog.ErrorContext(ctx, "Failed to handle message", slog.Any("error", err))

				return tgbotapi.NewMessage(chatID, "Sorry, I encountered an error while processing your request. Please try again later."), nil
			}
			return msgConfig, nil
		})
	}
}
