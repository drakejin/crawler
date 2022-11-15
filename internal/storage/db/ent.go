package db

import (
	"context"
	"database/sql"
	"os"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/drakejin/crawler/internal/storage/db/ent"
	"github.com/drakejin/crawler/internal/storage/db/ent/migrate"
)

type adapter struct {
	client *ent.Client
}

func New(db *sql.DB, debug bool) *adapter {
	drv := entsql.OpenDB("mysql", db)
	var options = []ent.Option{
		ent.Driver(drv),
	}
	if debug {
		entLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Str("dbClient", "entgo").
			Str("dialect", "mysql").
			Logger()

		options = append(
			options,
			ent.Driver(dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
				entLogger.Debug().Interface("action", i).Send()
			})),
		)
	}

	return &adapter{
		client: ent.NewClient(options...),
	}
}

func (a *adapter) Close() error {
	return a.client.Close()
}

func (a *adapter) Client() *ent.Client {
	return a.client
}

func (a *adapter) Migrate() error {
	start := time.Now()
	if err := a.client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(false),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Error().
			Err(err).
			Dur("duration", time.Since(start)).
			Msg("Error while in migration")
		return err
	}

	log.Info().Dur("duration", time.Since(start)).Send()
	return nil
}
