package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg CreateUserParams) (*User, error)
}

type ConduitStore struct {
	*Queries // implements Querier
	db       *pgxpool.Pool
}

func NewConduitStore(db *pgxpool.Pool) Store {
	return &ConduitStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *ConduitStore) CreateUserTx(
	ctx context.Context,
	arg CreateUserParams,
) (*User, error) {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	qtx := store.Queries.WithTx(tx)

	_, err = qtx.GetByPhoneNumber(ctx, arg.PhoneNumber)

	// 1. Real DB error? bail.
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 2. Found a user? bail with conflict.
	if err == nil {
		return nil, fmt.Errorf("phone number already registered")
	}

	newUser, err := qtx.CreateUser(ctx, CreateUserParams{
		Name:        arg.Name,
		PhoneNumber: arg.PhoneNumber,
		Password:    arg.Password,
		Birthdate:   arg.Birthdate,
	})

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err
	}

	return &newUser, nil

}
