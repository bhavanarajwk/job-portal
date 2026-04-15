package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"

	"job-portal/config"
)

// EmailMessage holds everything needed to send one email
type EmailMessage struct {
	To      []string
	Subject string
	Body    string // HTML
}

// SendEmail sends an HTML email via SMTP.
// It runs synchronously — call it inside a goroutine for non-blocking use.
func SendEmail(msg EmailMessage) error {
	cfg := config.AppConfig

	if cfg.SMTPUser == "" || cfg.SMTPPassword == "" {
		log.Println("[WARN] SMTP credentials not configured, skipping email")
		return nil
	}

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		cfg.SMTPFrom,
		strings.Join(msg.To, ", "),
		msg.Subject,
	)

	body := []byte(headers + msg.Body)
	addr := net.JoinHostPort(cfg.SMTPHost, cfg.SMTPPort)

	// Use STARTTLS (port 587)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}

	client, err := smtp.NewClient(conn, cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{ServerName: cfg.SMTPHost}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("starttls error: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	if err = client.Mail(cfg.SMTPUser); err != nil {
		return fmt.Errorf("smtp MAIL error: %w", err)
	}

	for _, to := range msg.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp RCPT error for %s: %w", to, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA error: %w", err)
	}

	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}

	return client.Quit()
}

// SendEmailAsync fires SendEmail in a goroutine so it never blocks the API
func SendEmailAsync(msg EmailMessage) {
	go func() {
		if err := SendEmail(msg); err != nil {
			log.Printf("[EMAIL ERROR] Failed to send to %v: %v", msg.To, err)
		} else {
			log.Printf("[EMAIL] Sent '%s' to %v", msg.Subject, msg.To)
		}
	}()
}
