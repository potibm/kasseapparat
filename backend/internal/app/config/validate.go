package config

import (
	"fmt"
	"log/slog"
	"net"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validDbFilename = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	validLocale     = regexp.MustCompile(`^[a-zA-Z]{2}-[A-Z]{2}$`)
	validCurrency   = regexp.MustCompile(`^[A-Z]{3}$`)
)

func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := c.App.Validate(); err != nil {
		return err
	}

	if err := c.Format.Validate(); err != nil {
		return err
	}

	if c.Jwt.Secret == DefaultJwtSecret || c.Jwt.Secret == "" {
		if c.App.Environment == "production" {
			return fmt.Errorf("JWT_SECRET is set to the default value, which is not allowed in production")
		} else {
			slog.Warn("JWT_SECRET is set to the default value. This is not recommended for production use.")
		}
	}

	return nil
}

func (f *AppConfig) Validate() error {
	if !validDbFilename.MatchString(f.DbFilename) {
		return fmt.Errorf("db_filename '%s' contains invalid characters", f.DbFilename)
	}

	if f.RedisURL != "" {
		if err := f.RedisURL.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisUrl) Validate() error {
	rString := string(*r)

	if !r.IsValid() {
		return fmt.Errorf("redis_url '%s' is not a valid URL", rString)
	}

	redisUrl := r.UrlObject()
	if redisUrl.Scheme != "redis" && redisUrl.Scheme != "rediss" {
		return fmt.Errorf(
			"redis_url '%s' has invalid scheme '%s' (expected 'redis' or 'rediss')",
			rString,
			redisUrl.Scheme,
		)
	}

	host, _, err := net.SplitHostPort(redisUrl.Host)
	if err != nil || host == "" {
		return fmt.Errorf("redis_url '%s' has missing host", rString)
	}

	return nil
}

func (f *FormatConfig) Validate() error {
	if err := f.Currency.Validate(); err != nil {
		return err
	}

	if err := f.Date.Validate(); err != nil {
		return err
	}

	return nil
}

func (f *DateFormatConfig) Validate() error {
	if !validLocale.MatchString(f.Locale) {
		return fmt.Errorf("date.locale '%s' is not a valid locale", f.Locale)
	}

	return nil
}

func (f *CurrencyFormatConfig) Validate() error {
	if !validLocale.MatchString(f.Locale) {
		return fmt.Errorf("currency.locale '%s' is not a valid locale", f.Locale)
	}

	if !validCurrency.MatchString(f.Code) {
		return fmt.Errorf("currency.code '%s' is not a valid ISO 4217 code", f.Code)
	}

	return nil
}
