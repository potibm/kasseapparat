package config

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/appleboy/gin-jwt/v3/store"
)

func (vr *VatRatesConfig) Json() string {
	jsonData, err := json.Marshal(*vr)
	if err != nil {
		return "[]"
	}

	return string(jsonData)
}

func (vr *DateFormatOptionsConfig) Json() string {
	jsonData, err := json.Marshal(*vr)
	if err != nil {
		return "[]"
	}

	return string(jsonData)
}

func (u *RedisUrl) UrlObject() *url.URL {
	if u == nil {
		return nil
	}

	parsedUrl, err := url.Parse(string(*u))
	if err != nil {
		return nil
	}

	return parsedUrl
}

func (u *RedisUrl) IsValid() bool {
	if u == nil {
		return true
	}

	_, err := url.ParseRequestURI(string(*u))

	return err == nil
}

func (r RedisUrl) JwtConfig() store.RedisConfig {
	url := r.UrlObject()

	if url == nil {
		return store.RedisConfig{}
	}

	path := url.Path
	if path == "" {
		path = "/0"
	}

	db, err := strconv.Atoi(path[1:])
	if err != nil {
		db = 0
	}

	password, _ := url.User.Password()

	return store.RedisConfig{
		Addr:     url.Host,
		Password: password,
		DB:       db,
	}
}
