package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	ReadPool  *pgxpool.Pool
	WritePool *pgxpool.Pool

	writerUri string
	readerUri string

	maxOpenConnections    int
	maxIdleConnections    int
	maxConnectionLifeTime time.Duration
	maxConnectionIdleTime time.Duration
}

func New(writerUri, readerUri string, opts ...Option) *Postgres {
	pg := &Postgres{
		writerUri:             writerUri,
		readerUri:             readerUri,
		maxOpenConnections:    defaultMaxOpenConnections,
		maxIdleConnections:    defaultMaxIdleConnections,
		maxConnectionLifeTime: defaultMaxConnectionLifetime,
		maxConnectionIdleTime: defaultMaxConnectionIdleTime,
	}

	for _, opt := range opts {
		opt(pg)
	}

	return pg
}

func (pg *Postgres) Connect(ctx context.Context) error {
	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	writeConfig, err := pgxpool.ParseConfig(pg.writerUri)
	if err != nil {
		return fmt.Errorf("failed to create write config: %w", err)
	}

	readConfig, err := pgxpool.ParseConfig(pg.readerUri)
	if err != nil {
		return fmt.Errorf("failed to create read config: %w", err)
	}

	writeConfig.MinConns = int32(pg.maxIdleConnections)
	readConfig.MinConns = int32(pg.maxIdleConnections)
	writeConfig.MaxConns = int32(pg.maxOpenConnections)
	readConfig.MaxConns = int32(pg.maxOpenConnections)

	writeConfig.MaxConnLifetime = pg.maxConnectionLifeTime
	readConfig.MaxConnLifetime = pg.maxConnectionLifeTime
	writeConfig.MaxConnLifetimeJitter = time.Duration(0.2 * float64(pg.maxConnectionLifeTime))
	readConfig.MaxConnLifetimeJitter = time.Duration(0.2 * float64(pg.maxConnectionLifeTime))

	writeConfig.MaxConnIdleTime = pg.maxConnectionIdleTime
	readConfig.MaxConnIdleTime = pg.maxConnectionIdleTime

	writePool, err := pgxpool.NewWithConfig(initCtx, writeConfig)
	if err != nil {
		return fmt.Errorf("failed to create write pool: %w", err)
	}

	readPool, err := pgxpool.NewWithConfig(initCtx, readConfig)
	if err != nil {
		writePool.Close()
		return fmt.Errorf("failed to create read pool: %w", err)
	}

	retryPolicy := backoff.NewExponentialBackOff()
	retryPolicy.MaxElapsedTime = 1 * time.Minute

	// confirm connectivity
	err = backoff.Retry(func() error {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer pingCancel()

		if err := writePool.Ping(pingCtx); err != nil {
			return fmt.Errorf("write pool ping failed: %w", err)
		}
		if err := readPool.Ping(pingCtx); err != nil {
			return fmt.Errorf("read pool ping failed: %w", err)
		}
		return nil
	}, retryPolicy)

	if err != nil {
		writePool.Close()
		readPool.Close()
		return fmt.Errorf("pinging pools failed: %w", err)
	}

	pg.WritePool = writePool
	pg.ReadPool = readPool

	return nil
}

func (p *Postgres) Close() error {
	p.ReadPool.Close()
	p.WritePool.Close()
	return nil
}
