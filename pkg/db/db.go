package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/repository/queries"
)

type DB struct {
	Config     *pgxpool.Config
	Pool       *pgxpool.Pool
	rawQueries *queries.Queries
}

func NewDB(cfg config.Config) *DB {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Panic().Msg("Failed to parse database config")
	}

	dbConfig.MaxConns = cfg.Database.Connections
	dbConfig.MinConns = cfg.Database.Connections

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Panic().Msg("Failed to connect to database")
	}

	return &DB{
		Config:     dbConfig,
		Pool:       pool,
		rawQueries: queries.New(pool),
	}
}

func (db *DB) Queries(ctx context.Context) *queries.Queries {
	if ctx.Value("tx") != nil {
		return queries.New(db.Pool).WithTx(ctx.Value("tx").(pgx.Tx))
	}
	return db.rawQueries
}

func (db *DB) WithTx(ctx context.Context, fn func(context.Context) error) error {

	if ctx.Value("tx") != nil {
		return errors.New("transaction already in progress")
	}
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	newCtx := context.WithValue(ctx, "tx", tx)
	defer tx.Rollback(newCtx)

	if err := fn(newCtx); err != nil {
		return err
	}

	return tx.Commit(newCtx)
}
