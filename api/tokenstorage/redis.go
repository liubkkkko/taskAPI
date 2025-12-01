package tokenstorage

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

type SessionMeta struct {
    UserID    int    `json:"user_id"`
    IP        string `json:"ip,omitempty"`
    UserAgent string `json:"user_agent,omitempty"`
    CreatedAt int64  `json:"created_at"`
    Signed    string `json:"signed,omitempty"`
}

func RedisStart(RedisAddr, RedisPassword, RedisDb string) {
    // existing initialization kept as before (DB parsing moved out)
    // existing code in your project sets DB from env; keep that
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     RedisAddr,
        Password: RedisPassword,
        DB:       0, // default, previous code parsed env; keep zero if env parsing elsewhere
    })
    // Check connect to Redis
    pong, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %s", err)
    }
    fmt.Printf("Connected to Redis: %s\n", pong)
}

// SaveSession зберігає метадані session за ключем "refresh:{jti}" та ставить TTL
func SaveSession(client *redis.Client, jti string, meta SessionMeta, ttl time.Duration) error {
    ctx := context.Background()
    key := fmt.Sprintf("refresh:%s", jti)
    data, _ := json.Marshal(meta)
    if err := client.Set(ctx, key, data, ttl).Err(); err != nil {
        return err
    }
    return nil
}

// GetSession повертає SessionMeta або помилку (redis.Nil якщо нема)
func GetSession(client *redis.Client, jti string) (SessionMeta, error) {
    ctx := context.Background()
    key := fmt.Sprintf("refresh:%s", jti)
    res, err := client.Get(ctx, key).Result()
    if err != nil {
        return SessionMeta{}, err
    }
    var meta SessionMeta
    if err := json.Unmarshal([]byte(res), &meta); err != nil {
        return SessionMeta{}, err
    }
    return meta, nil
}

// DeleteSession видаляє ключ
func DeleteSession(client *redis.Client, jti string) error {
    ctx := context.Background()
    key := fmt.Sprintf("refresh:%s", jti)
    return client.Del(ctx, key).Err()
}

// SessionExists перевіряє наявність
func SessionExists(client *redis.Client, jti string) (bool, error) {
    ctx := context.Background()
    key := fmt.Sprintf("refresh:%s", jti)
    n, err := client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return n == 1, nil
}

// --- Сумісні з існуючою логікою helper-и (залишені для backward-compat) ---

func CheckKeyExists(client *redis.Client, key string) (bool, error) {
    ctx := context.Background()
    result, err := client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}

func CheckValueExists(client *redis.Client, key, value string) (bool, error) {
    ctx := context.Background()
    val, err := client.Get(ctx, value).Result()
    if err != nil {
        return false, err
    }
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