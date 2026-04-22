package email

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/resend/resend-go/v2"
)

func resendAPIKey() (string, error) {
	key := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	if key == "" {
		return "", fmt.Errorf("RESEND_API_KEY is not set")
	}
	return key, nil
}

func attachmentsFromPaths(paths []string) []*resend.Attachment {
	var atts []*resend.Attachment
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			log.Printf("email attachment read %s: %v", p, err)
			continue
		}
		atts = append(atts, &resend.Attachment{
			Filename: filepath.Base(p),
			Content:  b,
		})
	}
	return atts
}

// sendHTMLResend sends HTML mail via Resend.
// replyToEmail, when non-empty, sets Reply-To so the recipient's mail client replies to that address.
func sendHTMLResend(toEmail, toName, subject, htmlBody, replyToEmail string) (string, error) {
	key, err := resendAPIKey()
	if err != nil {
		return "", err
	}
	replyTo := strings.TrimSpace(replyToEmail)
	client := resend.NewClient(key)
	req := &resend.SendEmailRequest{
		From:    resendFromAddress(),
		To:      []string{formatMailbox(toEmail, toName)},
		Subject: subject,
		Html:    htmlBody,
	}
	if replyTo != "" {
		req.ReplyTo = replyTo
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func formatMailbox(email, displayName string) string {
	email = strings.TrimSpace(email)
	name := strings.TrimSpace(displayName)
	if email == "" {
		return ""
	}
	if name == "" || strings.EqualFold(name, email) {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

func resendFromAddress() string {
	if f := strings.TrimSpace(os.Getenv("RESEND_FROM")); f != "" {
		return f
	}
	return "Redbud Way <onboarding@resend.dev>"
}

// SendProviderQuoteRequest emails the provider quote-request notification with optional image attachments.
func SendProviderQuoteRequest(providerEmail, providerDisplayName, subject, htmlBody string, attachmentPaths []string) error {
	key, err := resendAPIKey()
	if err != nil {
		return err
	}
	client := resend.NewClient(key)
	_, err = client.Emails.Send(&resend.SendEmailRequest{
		From:        resendFromAddress(),
		To:          []string{formatMailbox(providerEmail, providerDisplayName)},
		Subject:     subject,
		Html:        htmlBody,
		Attachments: attachmentsFromPaths(attachmentPaths),
	})
	return err
}

// sendTextResendWithAttachments sends a plain-text message (e.g. customer message to provider) via Resend.
func sendTextResendWithAttachments(toEmail, toDisplayName, subject, textBody, replyTo string, attachmentPaths []string) error {
	key, err := resendAPIKey()
	if err != nil {
		return err
	}
	client := resend.NewClient(key)
	req := &resend.SendEmailRequest{
		From:        resendFromAddress(),
		To:          []string{formatMailbox(toEmail, toDisplayName)},
		Subject:     subject,
		Text:        textBody,
		Attachments: attachmentsFromPaths(attachmentPaths),
	}
	if rt := strings.TrimSpace(replyTo); rt != "" {
		req.ReplyTo = rt
	}
	_, err = client.Emails.Send(req)
	return err
}
