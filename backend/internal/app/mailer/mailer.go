package mailer

import (
	"log/slog"
	"net/smtp"
	"net/url"
	"strconv"
)

const (
	defaultSmtpPort        = 587
	defaultSmtpUsername    = ""
	defaultSmtpPassword    = ""
	defaultFrom            = "kasseapparat@example.com"
	defaultSubjectPrefix   = "[Kasseapparat] "
	defaultFrontendBaseUrl = "http://localhost:3000"
)

type SmtpConfig struct {
	user     string
	password string
	host     string
	port     int
}

type Mailer struct {
	smtpConfig SmtpConfig

	disabled bool

	from          string
	subjectPrefix string

	frontendBaseUrl string
}

const (
	footerTemplate = "mail/_footer.txt"
)

func NewMailer(dsn string) (*Mailer, error) {
	frontendBaseUrl := defaultFrontendBaseUrl

	smtpConfig, err := SmtpConfigFromDsn(dsn)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		smtpConfig:      *smtpConfig,
		from:            defaultFrom,
		subjectPrefix:   defaultSubjectPrefix,
		frontendBaseUrl: frontendBaseUrl,
		disabled:        false,
	}, nil
}

func SmtpConfigFromDsn(dsn string) (*SmtpConfig, error) {
	user := defaultSmtpUsername
	password := defaultSmtpPassword
	port := defaultSmtpPort

	u, err := url.Parse(dsn)
	if err != nil {
		slog.Error("Error parsing Mail DSN", "error", err)
		println(err)

		return nil, err
	}

	host := u.Hostname()
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}

	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			slog.Error("Invalid port in Mail DSN", "error", err)

			return nil, err
		}
	}

	return &SmtpConfig{
		user:     user,
		password: password,
		host:     host,
		port:     port,
	}, nil
}

func (m *Mailer) SetFrom(from string) {
	m.from = from
}

func (m *Mailer) SetSubjectPrefix(prefix string) {
	m.subjectPrefix = prefix
}

func (m *Mailer) SetFrontendBaseUrl(url string) {
	m.frontendBaseUrl = url
}

func (m *Mailer) SetDisabled(disabled bool) {
	m.disabled = disabled
}

func (m *Mailer) SendMail(to string, subject string, body string) error {
	if m.disabled {
		slog.Info("Mailer is disabled, not sending email")

		return nil
	}

	header := "From: " + m.from + "\r\n" +
		"Subject: " + m.subjectPrefix + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"

	message := []byte(header + body)

	var auth smtp.Auth = nil
	if m.smtpConfig.user != defaultSmtpUsername && m.smtpConfig.password != defaultSmtpPassword {
		auth = smtp.PlainAuth("", m.smtpConfig.user, m.smtpConfig.password, m.smtpConfig.host)
	}

	err := smtp.SendMail(m.address(), auth, m.from, []string{to}, message)
	if err != nil {
		slog.Error("Error sending mail", "error", err)
	}

	return err
}

func (m *Mailer) address() string {
	return m.smtpConfig.host + ":" + strconv.Itoa(m.smtpConfig.port)
}
