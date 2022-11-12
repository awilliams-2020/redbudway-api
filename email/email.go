package email

import (
	_ "embed"
	"os"
	"strings"

	"github.com/go-gomail/gomail"
)

//go:embed html/password.html
var password string

func email(email, name, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "help@redbudway.com")
	m.SetAddressHeader("To", email, name)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("mail.redbudway.com", 25, "help@redbudway.com", "MerCedEsAmgGt22$")

	return d.DialAndSend(m)
}

func ForgotPassword(userEmail, name, token, accountType, userID string) error {
	body := password

	body = strings.Replace(body, "{SUB_DOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	body = strings.Replace(body, "{TOKEN}", token, -1)
	body = strings.Replace(body, "{ACCOUNT_TYPE}", accountType, -1)
	body = strings.Replace(body, "{USER_ID}", userID, -1)

	return email(userEmail, name, "Reset Password", body)
}
