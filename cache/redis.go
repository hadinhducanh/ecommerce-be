package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ecommerce-be/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// ConnectRedis kết nối đến Redis
func ConnectRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDB,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = rdb
	log.Println("✅ Redis connected successfully!")
	return nil
}

// CloseRedis đóng kết nối Redis
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// Set lưu giá trị vào cache với TTL
func Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = RedisClient.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get lấy giá trị từ cache
func Get(key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("failed to get cache: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete xóa key khỏi cache
func Delete(key string) error {
	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// DeletePattern xóa tất cả keys theo pattern
func DeletePattern(pattern string) error {
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := RedisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("Failed to delete key %s: %v", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to delete pattern: %w", err)
	}
	return nil
}

// Exists kiểm tra key có tồn tại không
func Exists(key string) (bool, error) {
	count, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

// Cache keys patterns
const (
	CategoryKeyPrefix    = "category:"
	CategoryListKey      = "categories:list"
	CategorySearchKey    = "categories:search:"
	UserKeyPrefix        = "user:"
	UserProfileKeyPrefix = "user:profile:"
	ProductKeyPrefix     = "product:"
	ProductListKey       = "products:list"
	ProductSearchKey     = "products:search:"
)

// Helper functions để tạo cache keys
func CategoryKey(id uint) string {
	return fmt.Sprintf("%s%d", CategoryKeyPrefix, id)
}

func UserKey(id uint) string {
	return fmt.Sprintf("%s%d", UserKeyPrefix, id)
}

func UserProfileKey(id uint) string {
	return fmt.Sprintf("%s%d", UserProfileKeyPrefix, id)
}

func ProductKey(id uint) string {
	return fmt.Sprintf("%s%d", ProductKeyPrefix, id)
}
