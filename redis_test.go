package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
)

const (
	REDIS_DATABASE = 0
)

func TestRedisListToHash(t *testing.T) {
	const (
		REDIS_KEY_DESTINATION_HASH_KEY = "jobs:servers:history:temp:active"
	)

	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")

	sourceListKeys := []string{
		"jobs:servers:history:*:processing",
		"jobs:servers:history:*:paused",
	}

	var (
		ctx          = context.Background()
		m            = map[string]string{}
		hashItemKeys = []string{}
	)

	db := redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       REDIS_DATABASE,
		Password: password,
	})

	for _, item := range sourceListKeys {
		hashItemKeys = append(hashItemKeys, scanKeys(ctx, db, item)...)
	}

	for _, key := range hashItemKeys {
		data, err := db.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			t.Logf("Failed to read list key %s: %v\n", key, err)
			continue
		}

		for _, value := range data {
			m[value] = "history"
		}
	}

	if err := db.HSet(ctx, REDIS_KEY_DESTINATION_HASH_KEY, m).Err(); err != nil {
		t.Fatalf("Failed to write hash key %s: %v\n", REDIS_KEY_DESTINATION_HASH_KEY, err)
	}
}

func scanKeys(ctx context.Context, db *redis.Client, pattern string) []string {
	var keys []string
	var cursor uint64

	for {
		var err error
		var matchedKeys []string

		matchedKeys, cursor, err = db.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			fmt.Printf("Failed to scan keys with pattern %s: %v\n", pattern, err)
			break
		}

		keys = append(keys, matchedKeys...)
		if cursor == 0 {
			break
		}
	}

	return keys
}
