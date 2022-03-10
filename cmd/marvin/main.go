package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/maitesin/marvin/config"
	httpx "github.com/maitesin/marvin/internal/infra/http"
	sqlx "github.com/maitesin/marvin/internal/infra/sql"
	"github.com/maitesin/marvin/internal/infra/sql/migrations"
	"github.com/maitesin/marvin/pkg/tracking/correos"
	"github.com/maitesin/marvin/pkg/tracking/dhl"
	"github.com/upper/db/v4/adapter/postgresql"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	cfg := config.NewConfig()

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbConn.Close()

	pgConn, err := postgresql.New(dbConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pgConn.Close()

	dbDriver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	migrations.RegisterMigrationDriver(migrationsFS)
	migrations, err := migrate.NewWithDatabaseInstance("embed://migrations", "marvin", dbDriver)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = migrations.Up()
	if err != nil && err.Error() != "no change" {
		fmt.Println(err)
		return
	}

	deliveriesRepository := sqlx.NewDeliveriesRepository(pgConn)
	_ = deliveriesRepository

	go func() {
		err = http.ListenAndServe(
			strings.Join([]string{cfg.HTTP.Host, cfg.HTTP.Port}, ":"),
			httpx.DefaultRouter(),
		)
		if err != nil {
			fmt.Printf("Failed to start service: %s\n", err.Error())
		}
	}()

	correosTracker, err := correos.NewTracker(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	dhlTracker, err := dhl.NewTracker(http.DefaultClient)
	if err != nil {
		panic(err)
	}
	_ = dhlTracker

	//events, err := dhlTracker.Track("CO902088319DE")
	//if err != nil {
	//	panic(err)
	//}

	//for _, event := range events {
	//	fmt.Printf("%s\n%s\n\n", event.Timestamp, event.Information)
	//}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		builder := strings.Builder{}
		events, err := correosTracker.Track(update.Message.Text)
		if err != nil {
			panic(err)
		}

		for _, event := range events {
			builder.WriteString(fmt.Sprintf("%s\n%s\n\n", event.Timestamp, event.Information))
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
