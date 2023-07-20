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
    // Виконати перевірку наявності ключа в Redis
    result, err := client.Exists(ctx, key).Result()
    if err != nil {
        // Обробити помилку
        return false, err
    }

    // Перевірка результату
    return result == 1, nil
}

func CheckValueExists(client *redis.Client, key, value string) (bool, error) {
    ctx := context.Background()

    // Виконати пошук значення по ключу в Redis
    val, err := client.Get(ctx, key).Result()
    if err != nil {
        // Обробити помилку
        return false, err
    }
	fmt.Println("val", val)
	fmt.Println("value", value)
	fmt.Println("val == value", val == value)
    // Перевірити співпадіння значення з вказаним value
    return val == value, nil
}