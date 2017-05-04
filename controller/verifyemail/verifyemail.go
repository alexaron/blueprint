package verifyemail

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/blue-jay/core/router"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/lib/mail"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/user"
)

// Load the routes.
func Load() {
	router.Get("/verify", VerifyEmail, acl.DisallowAuth)
	router.Get("/awaiting_verification", Awaiting, acl.DisallowAuth)
	router.Post("/resend_code", ResendCode, acl.DisallowAuth)
}

// VerifyEmail handles user email verification.
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	// Any flash messages can only be shown on some page, reliably redirecting there.
	defer http.Redirect(w, r, "/login", http.StatusFound)

	// Code must contains something and correspond to an existing user.
	invCodeErr := errors.New("Invalid verification code")
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		c.FlashError(invCodeErr)
		return
	}
	u, noRows, err := user.ByCode(c.DB, code)
	if err != nil {
		if noRows {
			c.FlashError(invCodeErr)
		} else {
			c.FlashErrorGeneric(err)
		}
		return
	}

	// Verifying a surely existing user.
	if _, err = user.Verify(c.DB, u.ID); err != nil {
		c.FlashErrorGeneric(err)
		return
	}
	c.FlashSuccess("Your email has successfully been verified. You can login now.")
}

// Awaiting shows a messages that a user is registered and a letter is sent.
func Awaiting(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("verifyemail/awaiting_verification")
	v.Render(w, r)
}

// ResendCode resends verification code.
func ResendCode(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	email := r.PostFormValue("email")

	u, _, err := user.ByEmail(c.DB, email)
	// Log it in case of error. It can even be ErrNoRows and it's still suspicious in this endpoint.
	if err != nil {
		log.Println("[ResendCode] Failed to fetch a user by email: " + err.Error())
	}
	if err != nil || u.Verified {
		c.FlashError(errors.New("Invalid email"))
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := mail.SendVerification(c.Config, email, u.VerificationCode); err != nil {
		c.FlashErrorGeneric(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/awaiting_verification", http.StatusFound)
}
