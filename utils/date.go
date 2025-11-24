package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgDate(s string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}, nil
}
