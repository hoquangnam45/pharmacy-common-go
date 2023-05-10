package data

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient(host string, port int, username string, password string) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Username: username, // No username
			Password: password, // no password set
			DB:       0,        // use default DB
		}),
	}
}

func (s *Redis) Get(ctx context.Context, namespace string, key string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if namespace == "" {
		return s.client.Get(ctx, key).Result()
	}
	return s.client.Get(ctx, fmt.Sprintf("%s:%s", namespace, key)).Result()
}
