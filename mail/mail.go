package mail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type SmtpConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func Send(c *SmtpConfig, toAddress []string, subject, body string) error {
	d := gomail.NewDialer(c.Host, c.Port, c.UserName, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", c.UserName)
	m.SetHeader("To", toAddress...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return d.DialAndSend(m)
}
