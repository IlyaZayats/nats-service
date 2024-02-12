package streaming

import (
	"github.com/avast/retry-go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewNatsStreamingConnection(url string) (*stan.Conn, error) {
	var sc stan.Conn
	err := retry.Do(func() error {
		var err error
		sc, err = stan.Connect("nats-streaming", "subscriber", stan.NatsURL(url))
		return err
	},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			logrus.Debugf("Retrying request after error: %v", err.Error())
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable connect to nats-streaming")
	}
	logrus.Println("Nats-streaming connected!")
	return &sc, nil
}
