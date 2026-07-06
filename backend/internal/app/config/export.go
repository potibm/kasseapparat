package config

import (
	"strconv"

	"github.com/appleboy/gin-jwt/v3/store"
)

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
