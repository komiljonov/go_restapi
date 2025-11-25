package db

import (
	"context"
)

type Querier interface {
	AllUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, id int32) (User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)

	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}
