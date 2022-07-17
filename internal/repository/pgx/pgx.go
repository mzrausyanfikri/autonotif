package pgx

import (
	"context"

	"github.com/aimzeter/autonotif/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(cfg config.Postgresql) (*pgxpool.Pool, error) {
	poolcfg, err := pgxpool.ParseConfig(cfg.FormatURL())
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), poolcfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
