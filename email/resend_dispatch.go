package email

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/resend/resend-go/v2"
)

// sendHTMLResendOrSMTP sends HTML mail via Resend when RESEND_API_KEY is set; otherwise SMTP.
// Returns the Resend email ID (empty string when falling back to SMTP).
// replyToEmail, when non-empty, sets Reply-To so the recipient's mail client replies to that address.
func sendHTMLResendOrSMTP(toEmail, toName, subject, htmlBody, replyToEmail string) (string, error) {
	key := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	replyTo := strings.TrimSpace(replyToEmail)
	if key == "" {
		return "", emailWithOptionalReplyTo(toEmail, toName, subject, htmlBody, replyTo)
	}
	client := resend.NewClient(key)
	req := &resend.SendEmailRequest{
		From:    resendFromAddress(),
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlBody,
	}
	if replyTo != "" {
		req.ReplyTo = replyTo
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		log.Printf("Resend send failed (%s), falling back to SMTP: %v", subject, err)
		return "", emailWithOptionalReplyTo(toEmail, toName, subject, htmlBody, replyTo)
	}
	return resp.Id, nil
}

func resendFromAddress() string {
	if f := strings.TrimSpace(os.Getenv("RESEND_FROM")); f != "" {
		return f
	}
	if u := strings.TrimSpace(os.Getenv("SMTP_USER")); u != "" {
		return "Redbud Way <" + u + ">"
	}
	return "Redbud Way <onboarding@resend.dev>"
}

// SendProviderQuoteRequestResendOrSMTP sends the provider "quote request" notification with optional image attachments.
// When RESEND_API_KEY is unset, uses the legacy gomail path (same as before).
func SendProviderQuoteRequestResendOrSMTP(providerEmail, providerDisplayName, subject, htmlBody string, attachmentPaths []string) error {
	key := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	if key == "" {
		return sendProviderQuoteRequestSMTP(providerEmail, providerDisplayName, subject, htmlBody, attachmentPaths)
	}

	var atts []*resend.Attachment
	for _, p := range attachmentPaths {
		b, err := os.ReadFile(p)
		if err != nil {
			log.Printf("quote request attachment read %s: %v", p, err)
			continue
		}
		atts = append(atts, &resend.Attachment{
			Filename: filepath.Base(p),
			Content:  b,
		})
	}

	client := resend.NewClient(key)
	_, err := client.Emails.Send(&resend.SendEmailRequest{
		From:        resendFromAddress(),
		To:          []string{providerEmail},
		Subject:     subject,
		Html:        htmlBody,
		Attachments: atts,
	})
	if err != nil {
		log.Printf("Resend quote-request send failed, falling back to SMTP: %v", err)
		return sendProviderQuoteRequestSMTP(providerEmail, providerDisplayName, subject, htmlBody, attachmentPaths)
	}
	return nil
}

func sendProviderQuoteRequestSMTP(providerEmail, providerDisplayName, subject, htmlBody string, attachmentPaths []string) error {
	host, port, user, password, err := loadSMTPConfig()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "service@redbudway.com")
	m.SetAddressHeader("To", providerEmail, providerDisplayName)
	m.SetHeader("Subject", subject)
	for _, p := range attachmentPaths {
		m.Attach(p)
	}
	m.SetBody("text/html", htmlBody)
	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}
