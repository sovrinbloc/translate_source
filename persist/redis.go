package persist

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	c "translate_source/config"
)

var (
	RedisClient = new(redis.Client)
)

func InitRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			c.Env.Vars["REDIS_HOST"],
			c.Env.Vars["REDIS_PORT"]),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := RedisClient.Ping().Result()

	if err == nil {
		log.Println("Redis Initialized")
		if err != nil {
			panic(err)
		}
	} else {
		log.Fatal(pong, err)
	}
	return RedisClient
}
