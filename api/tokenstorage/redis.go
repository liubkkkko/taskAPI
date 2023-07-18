package tokenstorage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func RedisStart(RedisAddr, RedisPassword, RedisDb string) {
	RedisDbInt, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDbInt,
	})

	// Check connect to Redis
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)
}
