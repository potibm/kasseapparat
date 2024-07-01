package mailer

import "strconv"

func (mailer *Mailer) SendChangePasswordTokenMail (to string, userId uint, username string, token string) error {

	body := `Hello `+username+`,

You have requested to change your password. Please use this link to change it:

`+ mailer.frontendBaseUrl + `/change-password?token=` + token + `&userId=` + strconv.Itoa(int(userId)) + `

This link will expire in 15 minutes.

If you did not request this, please ignore this email.

Thank you,
Kasseapparat
---
This is an automated email, please do not reply to this email.`

	return mailer.SendMail(to, "Change your password", body); 
}

func (mailer *Mailer) SendNewUserTokenMail (to string, userId uint, username string, token string) error {

	body := `Hello `+username+`,

Your account for Kasseapparat has been created. Please use this link to set your password

`+ mailer.frontendBaseUrl + `/change-password?token=` + token + `&userId=` + strconv.Itoa(int(userId)) + `

This link will expire in 3 hours.

If you did not request this, please ignore this email.

Thank you,
Kasseapparat
---
This is an automated email, please do not reply to this email.`

	return mailer.SendMail(to, "Account created", body); 
}

