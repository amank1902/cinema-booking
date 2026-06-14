package redis

import (
	"context"
	"log"
	"os"

	goredis "github.com/redis/go-redis/v9"
)

// NewClient creates a client for the given address with no password.
func NewClient(addr string) *goredis.Client {
	return NewClientWithPassword(addr, "")
}

// NewClientWithPassword creates a Redis client for the given address and password.
func NewClientWithPassword(addr, password string) *goredis.Client {
	opts := &goredis.Options{Addr: addr}
	if password != "" {
		opts.Password = password
	}
	rdb := goredis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("connected to redis at %s", addr)
	return rdb
}

// NewClientFromEnv builds a Redis client from environment variables:
// REDIS_ADDR (default: localhost:6379) and REDIS_PASSWORD (optional).
func NewClientFromEnv() *goredis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	pwd := os.Getenv("REDIS_PASSWORD")
	return NewClientWithPassword(addr, pwd)
}
