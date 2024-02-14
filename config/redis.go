package config

import "github.com/go-redis/redis"

// method to initialize redis
func Init_redis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}
