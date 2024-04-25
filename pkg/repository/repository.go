package repository

import (
	"context"
	"tikkin/pkg/repository/queries"
)

type Repository interface {
	Queries(ctx context.Context) *queries.Queries
}
