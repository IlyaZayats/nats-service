package cache

import (
	"github.com/IlyaZayats/servord/internal/interfaces"
	"github.com/avast/retry-go"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MemcachedCache struct {
	client *memcache.Client
}

func (mc *MemcachedCache) Get(key string) ([]byte, error) {
	it, err := mc.client.Get(key)
	if err != nil {
		return nil, err
	}
	val := it.Value
	return val, err
}

func (mc *MemcachedCache) Set(key string, value []byte) error {
	err := mc.client.Set(&memcache.Item{Key: key, Value: value})
	return err
}

func NewMemcachedCache(url string) (interfaces.Cache, error) {
	mc := memcache.New(url)
	err := retry.Do(func() error {
		var err error
		err = mc.Set(&memcache.Item{Key: "connection", Value: []byte("test")})
		return err
	},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			logrus.Debugf("Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable connect to memcached")
	}
	logrus.Println("Memcached connected!")
	return &MemcachedCache{client: mc}, nil
}
