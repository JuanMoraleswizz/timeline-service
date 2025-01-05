package database

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Database interface {
	GetDb() *gorm.DB
}

type RedisDatabase interface {
	GetRedis() *redis.Client
}
