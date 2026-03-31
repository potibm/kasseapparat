package config

import (
	"fmt"
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
		return nil, err
	}

	if u.Scheme != "redis" {
		return nil, fmt.Errorf("invalid scheme: %s", u.Scheme)
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
