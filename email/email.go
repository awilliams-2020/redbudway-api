package email

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-gomail/gomail"
)

// loadSMTPConfig reads outbound mail settings from the environment (same contract as production SMTP).
func loadSMTPConfig() (host string, port int, user string, password string, err error) {
	password = os.Getenv("SMTP_PASSWORD")
	if password == "" {
		return "", 0, "", "", fmt.Errorf("SMTP_PASSWORD is not set")
	}
	host = os.Getenv("SMTP_HOST")
	if host == "" {
		host = "mail.redbudway.com"
	}
	port = 587
	if ps := os.Getenv("SMTP_PORT"); ps != "" {
		p, convErr := strconv.Atoi(ps)
		if convErr != nil {
			return "", 0, "", "", fmt.Errorf("invalid SMTP_PORT: %w", convErr)
		}
		port = p
	}
	user = os.Getenv("SMTP_USER")
	if user == "" {
		user = "service@redbudway.com"
	}
	return host, port, user, password, nil
}

// sendMailMessage sends a fully built gomail message using SMTP_* env (attachments, custom From, etc.).
func sendMailMessage(m *gomail.Message) error {
	host, port, user, password, err := loadSMTPConfig()
	if err != nil {
		return err
	}
	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}

//go:embed html/password.html
var password string

//go:embed html/password-updated.html
var passwordUpdated string

//go:embed html/email-updated.html
var emailUpdated string

func email(toEmail, name, subject, body string) error {
	host, port, user, password, err := loadSMTPConfig()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", user, "Redbud Way")
	m.SetAddressHeader("To", toEmail, name)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, user, password)
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
