package mailer

func (mailer *Mailer) SendNotificationOnArrival(to string, username string) error {

	body := `Hello,

We would like to inform you that ` + username + ` has arrived.

Best regards,
Kasseapparat
---
This is an automated email, please do not reply to this email.`

	return mailer.SendMail(to, "Guest has arrived", body)
}
