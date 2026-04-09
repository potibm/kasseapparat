package config

import "net/url"

const redacted = "***REDACTED***"

func (c Config) RedactConfigForDisplay() Config {
	result := c

	result.Jwt.Secret = redacted
	result.Sumup.APIKey = redacted
	result.Sentry.DSN = redacted

	result.App.RedisURL = RedisURL(redactURLPassword(string(c.App.RedisURL)))
	result.Mailer.DSN = redactURLPassword(c.Mailer.DSN)

	return result
}

func redactURLPassword(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if parsedURL.User != nil {
		parsedURL.User = url.UserPassword(parsedURL.User.Username(), redacted)

		return parsedURL.String()
	}

	return rawURL
}
