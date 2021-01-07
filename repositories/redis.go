package repositories

import (
	"context"
	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"

	"github.com/github.com/steevehook/account-api/logging"
)

// RedisSettings represents the redis client settings
type RedisSettings struct {
	URL      string
	Password string
}

// NewRedisDriver creates a new instance of redis client
func NewRedisDriver(settings RedisSettings) (*redis.Client, error) {
	logger := logging.Logger
	client := redis.NewClient(&redis.Options{
		Addr:     settings.URL,
		Password: settings.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("could not ping redis server", zap.Error(err))
		return nil, err
	}

	return client, nil
}
