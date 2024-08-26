package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/hrz8/got/config"
	"github.com/hrz8/got/database"
	"github.com/hrz8/got/internal/storage/postgres"
	"github.com/hrz8/got/pkg/logger"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	postgresMigrationDir = "migrations"
	postgresDialect      = "postgres"
	migrationsTable      = "migrations"
)

func MigrateUp() error {
	cfg := config.New()
	logger := logger.New(logger.LogLevel(cfg.LogLevel))

	logger.Info("running migrate up command")
	pg := postgres.New(
		cfg.DatabaseURL,
		cfg.DatabaseURLReader,
		postgres.MaxOpenConnections(20),
		postgres.MaxIdleConnections(1),
		postgres.MaxConnectionLifeTime(300),
		postgres.MaxConnectionIdleTime(60),
	)

	defer closeDB(pg)

	if err := pg.Connect(context.Background()); err != nil {
		return fmt.Errorf("cannot connect to database: %v", err.Error())
	}

	goose.SetTableName(migrationsTable)
	goose.SetDialect(postgresDialect)
	goose.SetBaseFS(database.MigrationsFS)

	pool := stdlib.OpenDBFromPool(pg.WritePool)
	if err := goose.Up(pool, postgresMigrationDir); err != nil {
		return fmt.Errorf("failed when perform migration up: %v", err.Error())
	}

	return nil
}

func MigrateDown() error {
	cfg := config.New()
	logger := logger.New(logger.LogLevel(cfg.LogLevel))

	logger.Info("running migrate down command")
	pg := postgres.New(
		cfg.DatabaseURL,
		cfg.DatabaseURLReader,
		postgres.MaxOpenConnections(20),
		postgres.MaxIdleConnections(1),
		postgres.MaxConnectionLifeTime(300),
		postgres.MaxConnectionIdleTime(60),
	)

	defer closeDB(pg)

	if err := pg.Connect(context.Background()); err != nil {
		return fmt.Errorf("cannot connect to database: %v", err.Error())
	}

	goose.SetTableName(migrationsTable)
	goose.SetDialect(postgresDialect)
	goose.SetBaseFS(database.MigrationsFS)

	pool := stdlib.OpenDBFromPool(pg.WritePool)
	if err := goose.Down(pool, postgresMigrationDir); err != nil {
		return fmt.Errorf("failed when perform migration down: %v", err.Error())
	}

	return nil
}

func closeDB(pg *postgres.Postgres) {
	if err := pg.Close(); err != nil {
		log.Printf("failed to close the database: %v\n", err.Error())
	}
}
