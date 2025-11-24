package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	createUser(ctx context.Context, arg createUserParams) (pgconn.CommandTag, error)
}
