package language

import (
	"github.com/go-redis/redis"
	"fmt"
	"log"
)

var (
	redisClient = new(redis.Client)
)
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})


	pong, err := redisClient.Ping().Result()

	if err == nil {
		log.Println("redis initialized")
		if err != nil {
			panic(err)
		}

		err := redisClient.Set("key", "value", 0).Err()
		if err != nil {
			panic(err)
		}


		val, err := redisClient.Get("key").Result()
		if err != nil {
			panic(err)
		}
		if got, expected := val, "value"; got != expected {
			log.Panicf("incorrect return value, got %s, expected %s\n", got, expected)
		}

		return

	} else {
		fmt.Println(pong, err)
	}
}

func cacheGet(key string) string {
	err := redisClient.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	return ""
}
