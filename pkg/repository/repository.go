package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/repository/queries"
)

type Repositories interface {
	Queries(ctx context.Context) *queries.Queries
}

type Repository struct {
	db     *db.DB
	config *config.Config
}

func NewRepository(db *db.DB, config *config.Config) Repository {
	return Repository{db: db, config: config}
}

func (r *Repository) Queries(ctx context.Context) *queries.Queries {
	if ctx.Value("tx") != nil {
		return queries.New(r.db.Pool).WithTx(ctx.Value("tx").(pgx.Tx))
	}
	return queries.New(r.db.Pool)
}
