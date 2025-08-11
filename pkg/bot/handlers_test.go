package bot

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupHandler(t *testing.T) {
	// Create a service with mocked dependencies
	mockTokenSvc := NewMockService(t)
	svc := &Bot{
		token: "test-token",
		tg:    NewMocktgClient(t),
		svc:   mockTokenSvc,
	}

	// Call setupHandler
	handler := svc.setupHandler()

	// Verify that the handler is not nil
	assert.NotNil(t, handler, "Handler should not be nil")
}

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		setupMocks func(mockTokenSvc *MockService)
		name       string
		command    string
		wantText   string
		chatID     int64
		userID     int64
		wantErr    bool
	}{
		{
			name:    "start command",
			command: "start",
			setupMocks: func(mockTokenSvc *MockService) {
			},
			chatID:   123,
			userID:   456,
			wantText: welcomeMessage,
			wantErr:  false,
		},
		{
			name:    "help command",
			command: "help",
			setupMocks: func(mockTokenSvc *MockService) {
				// No mocks needed for help command
			},
			chatID:   123,
			userID:   456,
			wantText: helpMessage,
			wantErr:  false,
		},
		{
			name:    "unknown command",
			command: "unknown",
			setupMocks: func(mockTokenSvc *MockService) {
				// No mocks needed for unknown command
			},
			chatID:   123,
			userID:   456,
			wantText: unknownCommandMessage,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a service with mocked dependencies
			mockTokenSvc := NewMockService(t)
			svc := &Bot{
				token: "test-token",
				tg:    NewMocktgClient(t),
				svc:   mockTokenSvc,
			}

			// Setup mocks
			tt.setupMocks(mockTokenSvc)

			// Create a message with the command
			msg := &tgbotapi.Message{
				Text: "/" + tt.command,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: len(tt.command) + 1,
					},
				},
				Chat: &tgbotapi.Chat{
					ID: tt.chatID,
				},
				From: &tgbotapi.User{
					ID: tt.userID,
				},
			}

			// Call handleCommand
			resp, err := svc.handleCommand(context.Background(), msg)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Check response
			if tt.wantText != "" {
				assert.Equal(t, tt.wantText, resp.Text)
			}

			// For new_token success, check that the response contains token information
			if tt.command == "new_token" && !tt.wantErr && tt.wantText == "" {
				assert.Contains(t, resp.Text, "Your New API Token")
				assert.Contains(t, resp.Text, "token123")
				assert.Contains(t, resp.Text, "Valid until:")
			}
		})
	}
}
