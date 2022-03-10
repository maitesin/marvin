package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	private *tgbotapi.BotAPI
}

func NewBot(cfg Config) (Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		fmt.Println(err)
		return Bot{}, err
	}

	bot.Debug = true

	return Bot{private: bot}, nil
}

func (b *Bot) GetUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.private.GetUpdatesChan(u)
}

func (b *Bot) Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	return b.private.Send(msg)
}
