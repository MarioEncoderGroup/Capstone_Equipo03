package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis configuration for MisViaticos
type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Database  int
	TLSConfig *tls.Config
	Reset     bool
}

// Storage represents Redis storage for Fiber limiter
type Storage struct {
	db     *redis.Client
	config Config
}

// New creates a new Redis storage instance for MisViaticos
func New(config ...Config) *Storage {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = Config{
			Host:     "localhost",
			Port:     6379,
			Database: 0,
		}
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:  cfg.Username,
		Password:  cfg.Password,
		DB:        cfg.Database,
		TLSConfig: cfg.TLSConfig,
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	storage := &Storage{
		db:     rdb,
		config: cfg,
	}

	// Reset storage if requested
	if cfg.Reset {
		storage.Reset()
	}

	return storage
}

// Get retrieves a value from Redis
func (s *Storage) Get(key string) ([]byte, error) {
	ctx := context.Background()
	val, err := s.db.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

// Set stores a value in Redis with expiration
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	ctx := context.Background()
	return s.db.Set(ctx, key, val, exp).Err()
}

// Delete removes a key from Redis
func (s *Storage) Delete(key string) error {
	ctx := context.Background()
	return s.db.Del(ctx, key).Err()
}

// Reset clears all keys from the Redis database
func (s *Storage) Reset() error {
	ctx := context.Background()
	return s.db.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// GetClient returns the underlying Redis client for custom operations
func (s *Storage) GetClient() *redis.Client {
	return s.db
}

// Increment increments a key value (useful for rate limiting)
func (s *Storage) Increment(key string, exp time.Duration) (int64, error) {
	ctx := context.Background()
	
	// Use a pipeline for atomic operations
	pipe := s.db.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, exp)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return incr.Val(), nil
}

// TTL returns the time to live for a key
func (s *Storage) TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	return s.db.TTL(ctx, key).Result()
}

// Storage implements Redis storage for rate limiting