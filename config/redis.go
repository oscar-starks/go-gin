package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password for development
		DB:       0,  // default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
		log.Printf("Falling back to in-memory storage")
		RedisClient = nil
	} else {
		log.Println("Redis connected successfully!")
	}
}

func GetRedisClient() *redis.Client {
	return RedisClient
}
