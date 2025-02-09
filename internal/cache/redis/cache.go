package redis

import (
	"context"
	"encoding/json"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/go-redis/redis"
	"time"
)

func NewCache(cl *redis.Client) *cache {
	return &cache{
		cl: cl,
	}
}

type cache struct {
	cl *redis.Client
}

func (c *cache) Close() error {
	return c.cl.Close()
}

func (c *cache) Set(ctx context.Context, key string, val string) error {
	b, err := json.Marshal(val)
	if err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	if err = c.cl.Set(key, b, time.Hour*24).Err(); err != nil {
		panic(err)
	}

	return nil
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	b, err := c.cl.Get(key).Bytes()
	if err != nil {
		return "", utils.NewError(err.Error(), utils.InternalErr)
	}

	var res string
	if err = json.Unmarshal(b, &res); err != nil {
		return "", utils.NewError(err.Error(), utils.InternalErr)
	}

	return res, nil
}
