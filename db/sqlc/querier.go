package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	AllUsers(ctx context.Context) ([]User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
}
