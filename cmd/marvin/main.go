package main

import (
	"fmt"
	"net/http"

	"github.com/maitesin/marvin/pkg/tracking/correos"
	"github.com/maitesin/marvin/pkg/tracking/dhl"
)

func main() {
	_, err := correos.NewTracker(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	tracker, err := dhl.NewTracker(http.DefaultClient)
	if err != nil {
		panic(err)
	}
	events, err := tracker.Track("CO902088319DE")
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("%s\n%s\n\n", event.Timestamp, event.Information)
	}

	//bot, err := tgbotapi.NewBotAPI("")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//bot.Debug = true
	//
	//log.Printf("Authorized on account %s", bot.Self.UserName)
	//
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//updates := bot.GetUpdatesChan(u)
	//
	//for update := range updates {
	//	if update.Message == nil { // ignore any non-Message Updates
	//		continue
	//	}
	//
	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//	builder := strings.Builder{}
	//	events, err := tracker.Track(update.Message.Text)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	for _, event := range events {
	//		builder.WriteString(fmt.Sprintf("%s\n%s\n\n", event.Timestamp, event.Information))
	//	}
	//
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	bot.Send(msg)
	//}
}
