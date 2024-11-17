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

func CheckKeyExists(client *redis.Client, key string) (bool, error) {
	ctx := context.Background()
	// Check the presence of the key in Redis
	result, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	// check result
	return result == 1, nil
}

func CheckValueExists(client *redis.Client, key, value string) (bool, error) {
	ctx := context.Background()
	// search by key
	val, err := client.Get(ctx, value).Result()
	if err != nil {
		return false, err
	}
	// check if the key and value are the same
	return val == key, nil
}

func IdIfYouHaveToken(client *redis.Client, token string) (string, error) {
	ctx := context.Background()

	id, err := client.Get(ctx, token).Result()
	if err != nil {
		return "0", err
	}
	return id, nil
}
