package infrastructure

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
)

const dsn = `pool_max_conns=30 host=localhost port=5432 user=usr password=postgres dbname=films sslmode=disable`

func InitDatabase() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	return pgxpool.ConnectConfig(context.Background(), config)
}
