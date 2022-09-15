package util

import (
	"RD-Clone-API/pkg/config"
	"net/smtp"
)

func SendMail(subject, body, to string) {
	c := config.LoadConfig()
	auth := config.GetSmtp(c)

	message := []byte(subject + "\\n" + body)
	err := smtp.SendMail(c.SMTP.Host+":"+c.SMTP.Port, auth, c.SMTP.From, []string{to}, message)

	for err != nil {
		err = smtp.SendMail(c.SMTP.Host+":"+c.SMTP.Port, auth, c.SMTP.From, []string{to}, message)
	}
}
