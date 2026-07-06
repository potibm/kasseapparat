package config

import (
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
)

type RedisURL string

func (ru *RedisURL) URLObject() *url.URL {
	if ru == nil {
		return nil
	}

	parsedURL, err := url.Parse(string(*ru))
	if err != nil {
		return nil
	}

	return parsedURL
}

func (ru RedisURL) RedisOptions() *redis.Options {
	u := ru.URLObject()
	if u == nil {
		return nil
	}

	opt, err := redis.ParseURL(u.String())
	if err != nil {
		return nil
	}

	return opt
}

func (ru *RedisURL) Validate() error {
	if ru == nil {
		return nil
	}

	rString := string(*ru)
	if rString == "" {
		return fmt.Errorf("redis_url is empty")
	}

	_, err := redis.ParseURL(rString)
	if err != nil {
		return fmt.Errorf("invalid redis_url '%s': %w", rString, err)
	}

	return nil
}

func (ru *RedisURL) IsValid() bool {
	return ru.Validate() == nil
}

func (ru RedisURL) Redacted() RedisURL {
	return RedisURL(redactURLPassword(string(ru)))
}
