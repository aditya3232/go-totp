package configs

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper, log *logrus.Logger) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetString("REDIS_ADDR"),
		DB:       config.GetInt("REDIS_DB"),
		Password: config.GetString("REDIS_PASSWORD"),
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	return redisClient
}
