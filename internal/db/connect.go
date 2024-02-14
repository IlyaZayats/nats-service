package db

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewPostgresPool(dbUrl string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		return nil, errors.Wrap(err, "unable create pgxpool")
	}

	err = retry.Do(func() error {
		err := pool.Ping(context.Background())
		return err
	},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			logrus.Debugf("Retrying request after error: %v", err)
		}),
	)

	if err != nil {
		return nil, errors.Wrap(err, "unable ping database")
	}

	logrus.Println("Database connected!")
	return pool, nil
}
