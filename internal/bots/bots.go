package bots

import (
	"bot21/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bots struct {
	Bot *tgbotapi.BotAPI
}

func New(cfg *config.TelegramConfig) (*Bots, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = cfg.Debug

	return &Bots{
		Bot: bot,
	}, nil
}
