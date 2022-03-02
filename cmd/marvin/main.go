package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	"github.com/maitesin/marvin/config"
	sqlx "github.com/maitesin/marvin/internal/infra/sql"
	"github.com/maitesin/marvin/pkg/tracking/correos"
	"github.com/maitesin/marvin/pkg/tracking/dhl"
	"github.com/upper/db/v4/adapter/postgresql"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func init() {
	source.Register("embed", &driver{})
}

type driver struct {
	httpfs.PartialDriver
}

func (d *driver) Open(path string) (source.Driver, error) {
	err := d.PartialDriver.Init(http.FS(migrationsFS), path[len("embed://"):])
	if err != nil {
		return nil, err
	}

	return d, nil
}

func main() {
	cfg := config.NewConfig()

	fmt.Println("Opening connection")
	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbConn.Close()

	fmt.Println("PG connection")
	pgConn, err := postgresql.New(dbConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pgConn.Close()

	fmt.Println("PG Driver")
	dbDriver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Setup DB migrations")
	migrations, err := migrate.NewWithDatabaseInstance("embed://migrations", "marvin", dbDriver)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Run DB migrations")
	err = migrations.Up()
	if err != nil && err.Error() != "no change" {
		fmt.Println(err)
		return
	}

	fmt.Println("Done")
	deliveriesRepository := sqlx.NewDeliveriesRepository(pgConn)
	_ = deliveriesRepository

	_, err = correos.NewTracker(http.DefaultClient)
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
