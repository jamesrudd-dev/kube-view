package database

import "github.com/go-redis/redis"

var client *redis.Client

func ConnectRedis() (*redis.Client, error) {
	// connect to redis and test connection
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRedis() error {
	return client.Close()
}
