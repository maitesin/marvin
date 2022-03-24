package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maitesin/marvin/internal/domain"
	"github.com/maitesin/marvin/internal/infra/sql"
	"github.com/maitesin/marvin/pkg/tracking"
)

type Bot struct {
	ctx        context.Context
	private    *tgbotapi.BotAPI
	deliveries *sql.DeliveriesRepository
}

func NewBot(ctx context.Context, cfg Config, deliveries *sql.DeliveriesRepository) (Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		fmt.Println(err)
		return Bot{}, err
	}

	bot.Debug = true

	return Bot{
		ctx:        ctx,
		private:    bot,
		deliveries: deliveries,
	}, nil
}

func (b *Bot) GetUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.private.GetUpdatesChan(u)
}

func (b *Bot) Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	return b.private.Send(msg)
}

func (b *Bot) Listen(tracker tracking.Tracker) error {
	for update := range b.GetUpdatesChannel() {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		builder := strings.Builder{}
		events, err := tracker.Track(update.Message.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, event := range events {
			builder.WriteString(fmt.Sprintf("%s\n%s\n\n", event.Timestamp, event.Information))
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
		//msg.ReplyToMessageID = update.Message.MessageID

		err = b.deliveries.Insert(b.ctx, domain.NewDelivery(update.Message.Text, events))
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = b.Send(msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}
