package redis

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"os"
	"strings"

	goredis "github.com/redis/go-redis/v9"
)

// NewClient creates a client for the given address with no password.
func NewClient(addr string) *goredis.Client {
	return NewClientWithPassword(addr, "")
}

// NewClientWithPassword creates a Redis client for the given address and password.
func NewClientWithPassword(addr, password string, tlsEnabled ...bool) *goredis.Client {
	opts := &goredis.Options{Addr: addr}
	if password != "" {
		opts.Password = password
	}
	if len(tlsEnabled) > 0 && tlsEnabled[0] {
		opts.TLSConfig = &tls.Config{}
	}
	rdb := goredis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("connected to redis at %s", addr)
	return rdb
}

// NewClientFromEnv builds a Redis client from environment variables.
// It supports either a full REDIS_URL (like the Upstash CLI URL) or
// REDIS_ADDR + REDIS_PASSWORD.
// Examples:
//
//	REDIS_URL=redis://default:passwd@host:6379
//	REDIS_ADDR=host:6379
//	REDIS_PASSWORD=passwd
func NewClientFromEnv() *goredis.Client {
	if addr, pwd, tlsOn, ok := redisConfigFromEnv(); ok {
		return NewClientWithPassword(addr, pwd, tlsOn)
	}

	return NewClientWithPassword("localhost:6379", "")
}

func redisConfigFromEnv() (addr string, password string, tlsOn bool, ok bool) {
	raw := firstNonEmpty(os.Getenv("REDIS_URL"), os.Getenv("REDIS_ADDR"))
	if raw == "" {
		return "", "", false, false
	}

	raw = extractRedisURI(raw)
	if raw == "" {
		return "", "", false, false
	}

	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", false, false
		}

		addr = u.Host
		if u.User != nil {
			password, _ = u.User.Password()
		}
		scheme := strings.ToLower(u.Scheme)
		tlsOn = scheme == "rediss" || strings.EqualFold(os.Getenv("REDIS_TLS"), "true") || os.Getenv("REDIS_TLS") == "1" || strings.Contains(strings.ToLower(os.Getenv("REDIS_URL")), "--tls") || strings.Contains(strings.ToLower(os.Getenv("REDIS_ADDR")), "--tls")
		if addr == "" {
			return "", "", false, false
		}
		return addr, password, tlsOn, true
	}

	addr = strings.TrimSpace(raw)
	password = os.Getenv("REDIS_PASSWORD")
	tlsEnv := strings.ToLower(os.Getenv("REDIS_TLS"))
	tlsOn = tlsEnv == "1" || tlsEnv == "true"
	return addr, password, tlsOn, addr != ""
}

func extractRedisURI(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	if idx := strings.Index(raw, "redis://"); idx != -1 {
		raw = raw[idx:]
	} else if idx := strings.Index(raw, "rediss://"); idx != -1 {
		raw = raw[idx:]
	}

	raw = strings.Trim(raw, "'\"`[]()<> ")
	if idx := strings.IndexAny(raw, " \t\r\n"); idx != -1 {
		raw = raw[:idx]
	}
	raw = strings.Trim(raw, "'\"`[]()<> ")
	return raw
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
