package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/appleboy/gin-jwt/v3/store"
)

type RedisConfig url.URL

func loadRedisConfig() (*RedisConfig, error) {
	redisURL := getEnv("REDIS_URL", "")

	if redisURL == "" {
		return nil, nil
	}

	parsedUrl, err := url.Parse(redisURL)

	u, err := (*RedisConfig)(parsedUrl), err
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_URL: %w", err)
	}

	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return nil, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil || host == "" {
		return nil, errors.New("missing host in REDIS_URL")
	}

	return u, nil
}

func (r RedisConfig) JwtConfig() store.RedisConfig {
	path := r.Path
	if path == "" {
		path = "/0"
	}

	db, err := strconv.Atoi(path[1:])
	if err != nil {
		db = 0
	}

	password, _ := r.User.Password()

	return store.RedisConfig{
		Addr:     r.Host,
		Password: password,
		DB:       db,
	}
}
