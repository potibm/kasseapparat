package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmtpConfigFromDsnWithValidDsn(t *testing.T) {
	smtpConfig, err := SMTPConfigFromDSN("smtp://user:password@localhost:587")
	assert.NoError(t, err)
	assert.Equal(t, "user", smtpConfig.user)
	assert.Equal(t, "password", smtpConfig.password)
	assert.Equal(t, "localhost", smtpConfig.host)
	assert.Equal(t, 587, smtpConfig.port)
}

func TestSmtpConfigFromDsnWithInvalidDsn(t *testing.T) {
	_, err := SMTPConfigFromDSN("töst://invalid-dsn")
	assert.Error(t, err)
}

func TestSmtpConfigFromDsnWithInvalidPort(t *testing.T) {
	_, err := SMTPConfigFromDSN("smtp://user:password@localhost:abc")
	assert.Error(t, err)
}

func TestSmtpConfigFromDsnWithoutUsernameAndPassword(t *testing.T) {
	smtpConfig, err := SMTPConfigFromDSN("smtp://localhost:25")
	assert.NoError(t, err)
	assert.Equal(t, "", smtpConfig.user)
	assert.Equal(t, "", smtpConfig.password)
	assert.Equal(t, "localhost", smtpConfig.host)
	assert.Equal(t, 25, smtpConfig.port)
}

func TestSmtpConfigFromDsnWithoutPort(t *testing.T) {
	smtpConfig, err := SMTPConfigFromDSN("smtp://localhost")
	assert.NoError(t, err)
	assert.Equal(t, "", smtpConfig.user)
	assert.Equal(t, "", smtpConfig.password)
	assert.Equal(t, "localhost", smtpConfig.host)
	assert.Equal(t, 587, smtpConfig.port)
}
