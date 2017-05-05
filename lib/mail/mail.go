package mail

import (
	"fmt"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/core/email"
	"github.com/blue-jay/core/server"
	"gopkg.in/gomail.v2"
)

// SendVerification sends email that a user needs to verify their email.
func SendVerification(c env.Info, email, code string) error {
	return Send(
		c.Email,
		email,
		"Email verification",
		"Click this link to verify your email: "+schemeHost(c.Server)+"/verify?code="+code,
	)
}

// Send sends plaintext email.
func Send(info email.Info, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", info.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		info.Hostname,
		info.Port,
		info.Username,
		info.Password,
	)
	return d.DialAndSend(m)
}

// schemeHost returns a canonical URL representation for external access.
func schemeHost(c server.Info) string {
	schemeHost := "http"
	if c.UseHTTPS {
		schemeHost += "s"
	}
	schemeHost += "://" + c.Hostname
	if c.UseHTTPS {
		if c.HTTPSPort != 443 {
			schemeHost += fmt.Sprintf(":%d", c.HTTPSPort)
		}
	} else {
		if c.HTTPPort != 80 {
			schemeHost += fmt.Sprintf(":%d", c.HTTPPort)
		}
	}
	return schemeHost
}
