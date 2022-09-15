package config

import "net/smtp"

func GetSmtp(c *Config) smtp.Auth {
	return smtp.PlainAuth("", c.SMTP.Username, c.SMTP.Password, c.SMTP.Host)
}
