package email

import (
	_ "embed"
	"os"
	"strings"
)

//go:embed html/password.html
var password string

//go:embed html/password-updated.html
var passwordUpdated string

//go:embed html/email-updated.html
var emailUpdated string

func email(toEmail, name, subject, body string) error {
	return emailWithOptionalReplyTo(toEmail, name, subject, body, "")
}

// emailWithOptionalReplyTo sends HTML mail via Resend; replyTo is optional Reply-To header.
func emailWithOptionalReplyTo(toEmail, name, subject, body, replyTo string) error {
	_, err := sendHTMLResend(toEmail, name, subject, body, replyTo)
	return err
}

func ForgotPassword(userEmail, name, token, accountType string) error {
	body := password

	body = strings.Replace(body, "{SUB_DOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	body = strings.Replace(body, "{TOKEN}", token, -1)
	body = strings.Replace(body, "{ACCOUNT_TYPE}", accountType, -1)

	return email(userEmail, name, "Reset Password", body)
}

func PasswordUpdated(userEmail, name string) error {
	body := passwordUpdated

	body = strings.Replace(body, "{SUB_DOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	return email(userEmail, name, "Password Updated", body)
}

func EmailUpdated(userEmail, name, token, tradespersonID string) error {
	body := emailUpdated

	body = strings.Replace(body, "{SUB_DOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	body = strings.Replace(body, "{TOKEN}", token, -1)
	body = strings.Replace(body, "{TRADESPERSON_ID}", tradespersonID, -1)
	body = strings.Replace(body, "{EMAIL}", userEmail, -1)
	return email(userEmail, name, "Email Updated", body)
}
