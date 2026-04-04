package config

import (
	"fmt"
	"net"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, c.App.DbFilename)
	if !matched {
		return fmt.Errorf("db_filename '%s' contains invalid characters", c.App.DbFilename)
	}

	localeRegexp := regexp.MustCompile(`^[a-zA-Z]{2}-[A-Z]{2}$`)
	if !localeRegexp.MatchString(c.Format.Currency.Locale) {
		return fmt.Errorf("currency.locale '%s' is not a valid locale (expected format: xx-XX)", c.Format.Currency.Locale)
	}

	if !localeRegexp.MatchString(c.Format.Date.Locale) {
		return fmt.Errorf("date.locale '%s' is not a valid locale (expected format: xx-XX)", c.Format.Date.Locale)
	}

	currencyRegexp := regexp.MustCompile(`^[A-Z]{3}$`)
	if !currencyRegexp.MatchString(c.Format.Currency.Code) {
		return fmt.Errorf("currency.code '%s' is not a valid ISO 4217 currency code (expected format: XXX)", c.Format.Currency.Code)
	}

	if c.App.RedisURL != "" {
		if !c.App.RedisURL.IsValid() {
			return fmt.Errorf("redis_url '%s' is not a valid URL", c.App.RedisURL)
		}

		redisUrl := c.App.RedisURL.UrlObject()
		if redisUrl.Scheme != "redis" && redisUrl.Scheme != "rediss" {
			return fmt.Errorf("redis_url '%s' has invalid scheme '%s' (expected 'redis' or 'rediss')", c.App.RedisURL, redisUrl.Scheme)
		}

		host, _, err := net.SplitHostPort(redisUrl.Host)
		if err != nil || host == "" {
			return fmt.Errorf("redis_url '%s' has missing host", c.App.RedisURL)
		}
	}

	return nil
}
