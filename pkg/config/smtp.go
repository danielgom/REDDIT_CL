package config

import "net/smtp"

// SMTPAuth returns authentication for our mail server.
func SMTPAuth(c *Config) smtp.Auth {
	return smtp.PlainAuth("", c.SMTP.Username, c.SMTP.Password, c.SMTP.Host)
}
