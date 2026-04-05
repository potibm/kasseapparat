package config

import "net/url"

const redacted = "***REDACTED***"

func (c Config) RedactConfigForDisplay() Config {
	result := c

	result.Jwt.Secret = redacted
	result.Sumup.ApiKey = redacted
	result.Sentry.DSN = redacted

	result.App.RedisURL = RedisUrl(redactUrlPassword(string(c.App.RedisURL)))
	result.Mailer.DSN = redactUrlPassword(c.Mailer.DSN)

	return result
}

func redactUrlPassword(rawUrl string) string {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}

	if parsedUrl.User != nil {
		parsedUrl.User = url.UserPassword(parsedUrl.User.Username(), redacted)

		return parsedUrl.String()
	}

	return rawUrl
}
