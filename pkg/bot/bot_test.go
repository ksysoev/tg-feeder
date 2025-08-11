package bot

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	tests := []struct {
		cfg     *Config
		name    string
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
		{
			name:    "empty token",
			cfg:     &Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.cfg, NewMockService(t))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandle(t *testing.T) {
	mockTokenSvc := NewMockService(t)
	svc := &Bot{
		token: "test-token",
		tg:    NewMocktgClient(t),
		svc:   mockTokenSvc,
	}

	tests := []struct {
		name       string
		message    *tgbotapi.Message
		setupMocks func()
		wantText   string
		wantErr    bool
	}{
		{
			name: "start command",
			message: &tgbotapi.Message{
				Text: "/start",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 6,
					},
				},
				Chat: &tgbotapi.Chat{
					ID: 123,
				},
				From: &tgbotapi.User{
					ID: 456,
				},
			},
			setupMocks: func() {
			},
			wantText: welcomeMessage,
			wantErr:  false,
		},
		{
			name: "help command",
			message: &tgbotapi.Message{
				Text: "/help",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 5,
					},
				},
				Chat: &tgbotapi.Chat{
					ID: 123,
				},
				From: &tgbotapi.User{
					ID: 456,
				},
			},
			setupMocks: func() {},
			wantText:   helpMessage,
			wantErr:    false,
		},
		{
			name: "unknown command",
			message: &tgbotapi.Message{
				Text: "/unknown",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 8,
					},
				},
				Chat: &tgbotapi.Chat{
					ID: 123,
				},
				From: &tgbotapi.User{
					ID: 456,
				},
			},
			setupMocks: func() {},
			wantText:   unknownCommandMessage,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenSvc.ExpectedCalls = nil

			tt.setupMocks()

			msg, err := svc.Handle(context.Background(), tt.message)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.wantText != "" {
				assert.Equal(t, tt.wantText, msg.Text)
			}
		})
	}
}

func TestProcessUpdate(t *testing.T) {
	mockTg := NewMocktgClient(t)
	mockTokenSvc := NewMockService(t)

	cfg := &Config{
		Token: "test-token",
	}

	svc := &Bot{
		token: cfg.Token,
		tg:    mockTg,
		svc:   mockTokenSvc,
	}

	svc.handler = svc

	tests := []struct {
		update     *tgbotapi.Update
		setupMocks func()
		name       string
	}{
		{
			name: "nil message",
			update: &tgbotapi.Update{
				Message: nil,
			},
			setupMocks: func() {},
		},
		{
			name: "valid message",
			update: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: "/start",
					Entities: []tgbotapi.MessageEntity{
						{
							Type:   "bot_command",
							Offset: 0,
							Length: 6,
						},
					},
					Chat: &tgbotapi.Chat{
						ID: 123,
					},
					From: &tgbotapi.User{
						ID: 456,
					},
				},
			},
			setupMocks: func() {
				mockTg.EXPECT().Send(mock.Anything).Return(tgbotapi.Message{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTg.ExpectedCalls = nil
			mockTokenSvc.ExpectedCalls = nil

			tt.setupMocks()

			svc.processUpdate(context.Background(), tt.update)
		})
	}
}
