package mailer

import (
	"log"
	"net/smtp"
	"net/url"
	"strconv"
)

type Mailer struct {
	user     string
	password string
	host     string
	port     int

	from          string
	subjectPrefix string

	frontendBaseUrl string
}

func NewMailer(dsn string) *Mailer {
	user := ""
	password := ""
	host := "localhost"
	port := 1025
	frontendBaseUrl := "http://localhost:3000"

	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalln("Error parsing Mail DSN:", err)
	} else {
		host = u.Hostname()
		user = u.User.Username()
		password, _ = u.User.Password()
		port, _ = strconv.Atoi(u.Port())
	}

	return &Mailer{
		user:            user,
		password:        password,
		host:            host,
		port:            port,
		from:            "kasseapparat@example.com",
		subjectPrefix:   "[Kasseapparat] ",
		frontendBaseUrl: frontendBaseUrl,
	}
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

func (m *Mailer) address() string {
	return m.host + ":" + strconv.Itoa(m.port)
}

func (m *Mailer) SendMail(to string, subject string, body string) error {

	header := "From: " + m.from + "\r\n" +
		"Subject: " + m.subjectPrefix + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"

	message := []byte(header + body)

	var auth smtp.Auth = nil
	if m.user != "" && m.password != "" {
		auth = smtp.PlainAuth("", m.user, m.password, m.host)
	}
	err := smtp.SendMail(m.address(), auth, m.from, []string{to}, message)
	if err != nil {
		log.Println("Error sending mail:", err)
	}
	return err
}
