package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background() // Context

// Declaring struct of the cache instance
type CacheInstance struct {
	RedisClient *redis.Client
}

var Cache CacheInstance // variable of instace of cache

func ConnectCache() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "cars-cache:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})

	// Test Redis connection
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)

	Cache = CacheInstance{RedisClient: redisClient}
}
