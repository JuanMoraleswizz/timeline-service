package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"uala.com/timeline-service/config"
)

type redisDatabase struct {
	redisDb *redis.Client
}

func NewRedisDatabase(conf *config.Config) *redisDatabase {
	ctx := context.Background()
	redistClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password: "",
		DB:       0,
	})
	pong, err := redistClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
	return &redisDatabase{
		redisDb: redistClient,
	}
}

func (r *redisDatabase) GetRedis() *redis.Client {
	return r.redisDb
}
