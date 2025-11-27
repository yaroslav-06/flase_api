package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	Client *redis.Client
	ctx    context.Context
}

func Connect(port string, ctx context.Context) (Repo, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "",
		DB:       0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return Repo{}, fmt.Errorf("Failed to connect to Redis: %w", err)
	}

	// if appName, err := rdb.Get(ctx, "app").Result(); err != nil {
	// 	return Repo{}, fmt.Errorf("Error while reading app name: %w", err)
	// } else if appName != "blr" {
	// 	return Repo{}, fmt.Errorf("Can't connect to correct database, database app name: %s", appName)
	// }

	fmt.Println("Connected to redis successfully")
	return Repo{Client: rdb, ctx: ctx}, nil
}
