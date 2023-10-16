package redis

import (
	"github.com/redis/go-redis/v9"
	"suglider-auth/configs"
	"log"
	"context"
	"time"
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
		log.Println("Can not connect to redis:", err)
		panic(err)
	}

	log.Println("Connected to Redis successfully, redis master return:", pong)
}

// Redis SET
func Set(key, value string, ttl time.Duration) {

	err := rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Println("Redis SET data have something problem:", err)
		return
	}
}

// Redis GET
func Get(key string) string {

	value, err := rdb.Get(ctx, key).Result()

	// Check whether key exist or not
	if err == redis.Nil {
		log.Printf("Key '%s' does not exist.", key)
		return ""
	} else if err != nil {
		log.Println("Redis GET data failed:", err)
		return ""
	} else {
		log.Println("key:", value)
		return value
    }
}

// Redis EXISTS
func Exists(key string) bool {

	exists, err := rdb.Exists(ctx, key).Result()

	if err != nil {
		log.Println("Checking whether key exist or not happen something wrong:", err)
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
		log.Printf("Delete key(%s) failed: %v\n", key, err)
		return
	}

	log.Printf("Delete key(%s) successfully\n", key)
}

// Close redis connection
func Close() {
	rdb.Close()
	return
}
