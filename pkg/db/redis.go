package db

import (
	"context"
	"log"

	"github.com/AliAstanov/olx_clone/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {

	options := &redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DatabaseName,
	}

	client := redis.NewClient(options)

	_, err := client.Ping(ctx).Result()
	if err != nil{
		log.Println("Failed to connecting redis client:",err)
		return nil, err
	}
	return client,nil
}
