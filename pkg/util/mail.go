// Package util contains some utilities for the API.
package util

import (
	"fmt"
	"net/smtp"
	"strings"

	"RD-Clone-API/pkg/config"
)

func sendMail(subject, body, to string) {
	c := config.LoadConfig()
	auth := config.SMTPAuth(c)

	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("From: %s\n", c.SMTP.From))
	b.WriteString(fmt.Sprintf("To: %s\n", to))
	b.WriteString(fmt.Sprintf("Subject: %s\n\n", subject))
	b.WriteString(fmt.Sprintf("%s\n", body))

	message := []byte(b.String())
	addr := c.SMTP.Host + ":" + c.SMTP.Port

	err := smtp.SendMail(addr, auth, c.SMTP.From, []string{to}, message)

	for err != nil {
		err = smtp.SendMail(addr, auth, c.SMTP.From, []string{to}, message)
	}
}

// SendRegistrationEmail sends an email based on the subject, body and recipient.
func SendRegistrationEmail(token, to string) {
	basicMailRegistrationBody := "Thank you for signing up to Spring reddit service," +
		" please click the link below to activate your account\n\n %s"
	sendMail("Activate Spring reddit CL account", fmt.Sprintf(basicMailRegistrationBody, token), to)
}
