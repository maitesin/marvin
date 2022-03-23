package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/maitesin/marvin/config"
	httpx "github.com/maitesin/marvin/internal/infra/http"
	"github.com/maitesin/marvin/internal/infra/pinger"
	sqlx "github.com/maitesin/marvin/internal/infra/sql"
	"github.com/maitesin/marvin/internal/infra/sql/migrations"
	"github.com/maitesin/marvin/internal/infra/telegram"
	"github.com/maitesin/marvin/pkg/tracking/correos"
	"github.com/maitesin/marvin/pkg/tracking/dhl"
	"github.com/upper/db/v4/adapter/postgresql"
	"net/http"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()

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
		pinger.NewPinger(cfg.Pinger.Address, cfg.Pinger.Frequency).Start(ctx)
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

	bot, err := telegram.NewBot(cfg.Telegram)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		err = bot.Listen(correosTracker)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	err = http.ListenAndServe(
		strings.Join([]string{cfg.HTTP.Host, cfg.HTTP.Port}, ":"),
		httpx.DefaultRouter(),
	)
	if err != nil {
		fmt.Printf("Failed to start service: %s\n", err.Error())
	}
}
