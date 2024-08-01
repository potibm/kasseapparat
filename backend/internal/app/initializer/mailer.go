package initializer

import (
	"os"

	"github.com/potibm/kasseapparat/internal/app/mailer"
)

func InitializeMailer() mailer.Mailer {
	mailDsn := os.Getenv("MAIL_DSN")
	if mailDsn == "" {
		mailDsn = "smtp://user:password@localhost:1025"
	}
	mailFrom := os.Getenv("MAIL_FROM")
	mailSubjectPrefix := os.Getenv("MAIL_SUBJECT_PREFIX")
	frontendBaseUrl := os.Getenv("FRONTEND_URL")

	mailer := mailer.NewMailer(mailDsn)

	if mailFrom != "" {
		mailer.SetFrom(mailFrom)
	}
	if mailSubjectPrefix != "" {
		mailer.SetSubjectPrefix(mailSubjectPrefix)
	}
	if frontendBaseUrl != "" {
		mailer.SetFrontendBaseUrl(frontendBaseUrl)
	}

	return *mailer
}
