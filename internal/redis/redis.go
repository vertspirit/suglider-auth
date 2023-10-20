package redis

import (
	"github.com/redis/go-redis/v9"
	"suglider-auth/configs"
	"log/slog"
	"context"
	"time"
	"fmt"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     configs.ApplicationConfig.Redis.Host + ":" + configs.ApplicationConfig.Redis.Port,
		Password: configs.ApplicationConfig.Redis.Password,
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		errorMessage := fmt.Sprintf("Can not connect to redis: %v", err)
		slog.Error(errorMessage)
		
		panic(err)
	}

	slog.Info(fmt.Sprintf("Connected to Redis successfully, redis master return: %s", pong))

}

// Redis SET
func Set(key, value string, ttl time.Duration) {

	err := rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		errorMessage := fmt.Sprintf("Redis SET data have something problem: %v", err)
		slog.Error(errorMessage)

		return
	}
}

// Redis GET
func Get(key string) string {

	value, err := rdb.Get(ctx, key).Result()

	// Check whether key exist or not
	if err == redis.Nil {
		slog.Info(fmt.Sprintf("Key '%s' does not exist.", key))
		return ""
	} else if err != nil {
		errorMessage := fmt.Sprintf("Redis GET data failed: %v", err)
		slog.Error(errorMessage)
		return ""
	} else {
		slog.Info(fmt.Sprintf("key: %s", value))
		return value
    }
}

// Redis EXISTS
func Exists(key string) bool {

	exists, err := rdb.Exists(ctx, key).Result()

	if err != nil {
		errorMessage := fmt.Sprintf("Checking whether key exist or not happen something wrong: %v", err)
		slog.Error(errorMessage)
	}

	// Check whether key exist or not
	if exists == 1 {
		return true
	} else {
		return false
	}
}

// Redis DELETE
func Delete(key string) {

	err := rdb.Del(ctx, key).Err()
	if err != nil {
		errorMessage := fmt.Sprintf("Delete key(%s) failed: %v", key, err)
		slog.Error(errorMessage)
		
		return
	}

	slog.Info(fmt.Sprintf("Delete key(%s) successfully.", key))
}

// Close redis connection
func Close() {
	rdb.Close()
	return
}
