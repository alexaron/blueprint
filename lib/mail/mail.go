package mail

// TODO: Use gopkg.in/gomail.v2 instead of the default email sender.

import (
	"fmt"
	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/core/server"
)

// SendVerification sends email that a user needs to verify his email.
func SendVerification(c env.Info, email, code string) error {
	return c.Email.Send(
		email,
		"Email verification",
		"Click this link to verify your email: "+schemeHost(c.Server)+"/verify?code="+code,
	)
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
