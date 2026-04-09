package config

import (
	"net/url"
	"strconv"

	"github.com/appleboy/gin-jwt/v3/store"
)

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

func (ru *RedisURL) IsValid() bool {
	if ru == nil {
		return true
	}

	_, err := url.ParseRequestURI(string(*ru))

	return err == nil
}

func (ru RedisURL) JwtConfig() store.RedisConfig {
	u := ru.URLObject()

	if u == nil {
		return store.RedisConfig{}
	}

	path := u.Path
	if path == "" {
		path = "/0"
	}

	db, err := strconv.Atoi(path[1:])
	if err != nil {
		db = 0
	}

	password, _ := u.User.Password()

	return store.RedisConfig{
		Addr:     u.Host,
		Password: password,
		DB:       db,
	}
}
