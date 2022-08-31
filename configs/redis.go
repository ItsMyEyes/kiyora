package configs

import (
	"context"
	"fmt"
	"myself_framwork/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

func init() {
	if utils.GetEnv("USE_REDIS", "OFF") == "ON" {
		connect()
	}
}

var (
	ClientRedis *redis.Client
	ContextTim  = context.Background()
)

// create connection to redis
func connect() *redis.Client {
	ClientRedis = redis.NewClient(&redis.Options{
		Addr:     utils.Getenv("REDIS_ADDR", "192.168.111.131:6379"),
		Password: utils.Getenv("REDIS_PASSWORD", ""), // no password set
		DB:       0,                                  // use default DB
	})

	pong, err := ClientRedis.Ping(ContextTim).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, "connected to rds")
	return ClientRedis
}

func SetKey(key string, i interface{}, ttl time.Duration) *redis.StringCmd {
	if ttl == 0 {
		ClientRedis.Set(ContextTim, key, i, 0)
	} else {
		ClientRedis.Set(ContextTim, key, i, ttl)
	}

	return ClientRedis.Get(ContextTim, key)
}

func GetKey(key string) *redis.StringCmd {
	return ClientRedis.Get(ContextTim, key)
}

func FindKey(key string) *redis.StringSliceCmd {
	return ClientRedis.Keys(ContextTim, key)
}

func RemoveKey(key string) bool {
	ClientRedis.Del(ContextTim, key)
	return true
}

func HasKey(key string) bool {
	return ClientRedis.Exists(ContextTim, key).Val() == 1
}
