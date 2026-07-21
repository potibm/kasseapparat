package config

import (
	"fmt"
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

	if err := c.Auth.Validate(); err != nil {
		return err
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

func (a *AuthConfig) Validate() error {
	if a.Mode == "proxy" {
		if a.ProxyHeader == "" {
			return fmt.Errorf("auth.proxy_header is required when mode is 'proxy'")
		}

		return nil
	}

	if a.Mode == "oidc" {
		return a.validateOIDC()
	}

	return nil
}

func (a *AuthConfig) validateOIDC() error {
	if a.OidcIssuer == "" {
		return fmt.Errorf("auth.oidc_issuer is required when mode is 'oidc'")
	}

	if a.OidcClientID == "" {
		return fmt.Errorf("auth.oidc_client_id is required when mode is 'oidc'")
	}

	if a.OidcClientSecret == "" {
		return fmt.Errorf("auth.oidc_client_secret is required when mode is 'oidc'")
	}

	if a.OidcCallbackURL == "" {
		return fmt.Errorf("auth.oidc_callback_url is required when mode is 'oidc'")
	}

	if len(a.SessionSecret) < MinSessionSecretLength {
		return fmt.Errorf(
			"auth.session_secret must be at least %d characters when mode is 'oidc'",
			MinSessionSecretLength,
		)
	}

	return nil
}
