package initializer

import (
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/mailer"
)

func InitializeMailer(mailerConfig config.MailerConfig) mailer.Mailer {
	mail, err := mailer.NewMailer(mailerConfig.DSN)
	if err != nil {
		panic(err)
	}

	if mailerConfig.FromEmail != "" {
		mail.SetFrom(mailerConfig.FromEmail)
	}

	if mailerConfig.MailSubjectPrefix != "" {
		mail.SetSubjectPrefix(mailerConfig.MailSubjectPrefix)
	}

	if mailerConfig.FrontendURL != "" {
		mail.SetFrontendBaseUrl(mailerConfig.FrontendURL)
	}

	return *mail
}
