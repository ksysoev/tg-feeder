package bot

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksysoev/make-it-public-tgbot/pkg/bot/middleware"
)

const (
	welcomeMessage = `Wel`
	helpMessage    = `Available Commands:

/start - Show welcome message
/help - Display this help message
`
	unknownCommandMessage = "❓ Unknown command.\n\nUse /help to see the list of available commands."
)

// Handler defines the interface for processing and responding to incoming messages in a Telegram bot context.
// It handles a message by performing necessary processing and returns the configuration for the outgoing message or an error.
// ctx is the context for managing request lifecycle and cancellation.
// message is the incoming Telegram message to be processed.
// Returns a configured message object for sending a response and an error if processing fails.
type Handler interface {
	Handle(ctx context.Context, message *tgbotapi.Message) (tgbotapi.MessageConfig, error)
}

// setupHandler initializes and configures the request handler with specified middleware components.
// It applies middleware for request reduction, concurrency throttling, metric collection, and error handling,
// ensuring proper management of requests and enhanced error messages.
// Returns a Handler that processes messages with the applied middleware stack.
func (s *Bot) setupHandler() Handler {
	h := middleware.Use(
		s,
		middleware.WithThrottler(30),
		middleware.WithRequestSequencer(),
		middleware.WithMetrics(),
		middleware.WithErrorHandling(),
	)

	return h
}

// Handle processes incoming telegram messages, handles commands, text messages, and generates appropriate responses.
func (s *Bot) Handle(ctx context.Context, msg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	slog.DebugContext(ctx, "Handling message", slog.Any("message", msg))
	if msg.Command() != "" {
		resp, err := s.handleCommand(ctx, msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("failed to handle command: %w", err)
		}

		return resp, nil
	}

	return tgbotapi.MessageConfig{}, fmt.Errorf("not implemented")
}

// handleCommand handles Telegram command messages and generates an appropriate response based on the command received.
func (s *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	switch msg.Command() {
	case "start":
		return newTextMessage(msg.Chat.ID, welcomeMessage), nil
	case "help":
		return newTextMessage(msg.Chat.ID, helpMessage), nil
	case "add":
		_, err := s.svc.AddFeed(ctx, msg.Text)

		return tgbotapi.MessageConfig{}, err
	default:
		return newTextMessage(msg.Chat.ID, unknownCommandMessage), nil
	}
}
