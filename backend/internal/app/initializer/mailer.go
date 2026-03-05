package initializer

import (
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/mailer"
)

func InitializeMailer(mailerConfig config.MailerConfig) mailer.Mailer {
	mailer, err := mailer.NewMailer(mailerConfig.DSN)
	if err != nil {
		panic(err)
	}

	if mailerConfig.FromEmail != "" {
		mailer.SetFrom(mailerConfig.FromEmail)
	}

	if mailerConfig.MailSubjectPrefix != "" {
		mailer.SetSubjectPrefix(mailerConfig.MailSubjectPrefix)
	}

	if mailerConfig.FrontendURL != "" {
		mailer.SetFrontendBaseUrl(mailerConfig.FrontendURL)
	}

	return *mailer
}
