package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/potibm/kasseapparat/templates"
)

const (
	changePasswordSubject = "Change your password"
	accountCreatedSubject = "Account created"
)

func generateChangePasswordLink(baseUrl string, token string, userId uint) string {
	return fmt.Sprintf("%s/change-password?token=%s&userId=%d", baseUrl, token, userId)
}

func (mailer *Mailer) SendChangePasswordTokenMail(to string, userId uint, username string, token string) error {
	return mailer.sendTokenMail(to, userId, username, token, "mail/token_change_password.txt", changePasswordSubject)
}

func (mailer *Mailer) SendNewUserTokenMail(to string, userId uint, username string, token string) error {

	return mailer.sendTokenMail(to, userId, username, token, "mail/token_new_user.txt", accountCreatedSubject)
}

func (mailer *Mailer) sendTokenMail(to string, userId uint, username string, token string, templateFilename string, subject string) error {

	template, err := template.ParseFS(
		templates.MailTemplateFiles,
		templateFilename,
		footerTemplate,
	)
	if err != nil {
		return fmt.Errorf("Failed to parse email template: %w", err)
	}

	data := map[string]interface{}{
		"Username": username,
		"Link":     generateChangePasswordLink(mailer.frontendBaseUrl, token, userId),
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		return fmt.Errorf("Failed to execute email template: %w", err)
	}

	return mailer.SendMail(to, subject, body.String())
}
