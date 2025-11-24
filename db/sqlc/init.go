package db

import (
	"context"
	"log"
	"restapi/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// func Connect(config config.Config) *pgx.Conn {
// 	conn, err := pgx.Connect(context.Background(), config.DBUrl)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database: %v\n", err)
// 	}
// 	return conn
// }

func Connect(config config.Config) *pgxpool.Pool {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(config.DBUrl)
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.HealthCheckPeriod = 30

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	return pool

}
