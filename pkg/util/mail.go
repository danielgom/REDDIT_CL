// Package util contains some utilities for the API.
package util

import (
	"net/smtp"

	"RD-Clone-API/pkg/config"
)

// SendMail sends an email based on the subject, body and recipient.
func SendMail(subject, body, to string) {
	c := config.LoadConfig()
	auth := config.SMTPAuth(c)

	message := []byte(subject + "\\n" + body)
	err := smtp.SendMail(c.SMTP.Host+":"+c.SMTP.Port, auth, c.SMTP.From, []string{to}, message)

	for err != nil {
		err = smtp.SendMail(c.SMTP.Host+":"+c.SMTP.Port, auth, c.SMTP.From, []string{to}, message)
	}
}
