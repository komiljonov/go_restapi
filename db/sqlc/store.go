package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg createUserParams) (pgconn.CommandTag, error)
}
