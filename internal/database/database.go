package database

import "github.com/go-redis/redis"

var client *redis.Client

func InitialConnectRedis(clusterDatabase int) (*redis.Client, error) {
	// connect to redis and test connection
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       clusterDatabase,
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

func ChangeRedisDatabase(client *redis.Client, db_num int) (*redis.Client, error) {
	client.Close()

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       db_num,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
