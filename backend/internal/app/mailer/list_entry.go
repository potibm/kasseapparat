package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/potibm/kasseapparat/templates"
)

const (
	arrivalNotificationSubject = "Guest has arrived 🔔"
)

func (mailer *Mailer) SendNotificationOnArrival(to, username string) error {
	tpl, err := template.ParseFS(
		templates.MailTemplateFiles,
		"mail/notification_on_arrival.txt",
		footerTemplate,
	)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer

	err = tpl.Execute(&body, map[string]string{
		"Username": username,
	})
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return mailer.SendMail(to, arrivalNotificationSubject, body.String())
}
