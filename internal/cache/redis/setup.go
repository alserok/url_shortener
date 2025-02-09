package redis

import "github.com/go-redis/redis"

func MustConnect(dsn string) *redis.Client {
	cl := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})

	if err := cl.Ping().Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	return cl
}
