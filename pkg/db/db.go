package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/repository/queries"
)

type DB struct {
	Config  *pgxpool.Config
	Pool    *pgxpool.Pool
	Queries *queries.Queries
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
		Config:  dbConfig,
		Pool:    pool,
		Queries: queries.New(pool),
	}
}
