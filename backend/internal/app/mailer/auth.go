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

func generateChangePasswordLink(baseURL, token string, userID int) string {
	return fmt.Sprintf("%s/change-password?token=%s&userId=%d", baseURL, token, userID)
}

func (mailer *Mailer) SendChangePasswordTokenMail(to string, userID int, username, token string) error {
	return mailer.sendTokenMail(to, userID, username, token, "mail/token_change_password.txt", changePasswordSubject)
}

func (mailer *Mailer) SendNewUserTokenMail(to string, userID int, username, token string) error {
	return mailer.sendTokenMail(to, userID, username, token, "mail/token_new_user.txt", accountCreatedSubject)
}

func (mailer *Mailer) sendTokenMail(
	to string,
	userID int,
	username string,
	token string,
	templateFilename string,
	subject string,
) error {
	tpl, err := template.ParseFS(
		templates.MailTemplateFiles,
		templateFilename,
		footerTemplate,
	)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	data := map[string]any{
		"Username": username,
		"Link":     generateChangePasswordLink(mailer.frontendBaseURL, token, userID),
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return mailer.SendMail(to, subject, body.String())
}
