package email

import (
	_ "embed"
	"os"
	"strings"

	"github.com/go-gomail/gomail"
)

//go:embed html/password.html
var password string

//go:embed html/password-updated.html
var passwordUpdated string

//go:embed html/email-updated.html
var emailUpdated string

func email(email, name, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "service@redbudway.com")
	m.SetAddressHeader("To", email, name)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("mail.redbudway.com", 587, "service@redbudway.com", "MerCedEsAmgGt22$")

	return d.DialAndSend(m)
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
