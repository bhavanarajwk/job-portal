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

// SendEmail sends an HTML email via SMTP (synchronous).
// Call via SendEmailAsync for non-blocking use.
func SendEmail(msg EmailMessage) error {
	// Snapshot config values at call time — safe for goroutines
	cfg := config.AppConfig
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}

	host     := cfg.SMTPHost
	port     := cfg.SMTPPort
	user     := cfg.SMTPUser
	password := cfg.SMTPPassword
	from     := cfg.SMTPFrom

	if user == "" || password == "" {
		log.Println("[WARN] SMTP credentials not configured, skipping email")
		return nil
	}

	auth := smtp.PlainAuth("", user, password, host)

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		from,
		strings.Join(msg.To, ", "),
		msg.Subject,
	)

	fullBody := []byte(headers + msg.Body)
	addr := net.JoinHostPort(host, port)

	// Dial TCP
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Close()

	// Upgrade to TLS via STARTTLS
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("starttls error: %w", err)
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	// Set sender
	if err = client.Mail(user); err != nil {
		return fmt.Errorf("smtp MAIL FROM error: %w", err)
	}

	// Set recipients
	for _, to := range msg.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp RCPT TO error for %s: %w", to, err)
		}
	}

	// Write body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA error: %w", err)
	}
	if _, err = w.Write(fullBody); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("smtp body close error: %w", err)
	}

	return client.Quit()
}

// SendEmailAsync fires SendEmail in a background goroutine.
// Errors are logged but never returned — the API response is never blocked.
func SendEmailAsync(msg EmailMessage) {
	// Capture a copy of the message for the goroutine
	m := msg
	go func() {
		if err := SendEmail(m); err != nil {
			log.Printf("[EMAIL ERROR] to=%v subject=%q err=%v", m.To, m.Subject, err)
		} else {
			log.Printf("[EMAIL OK] to=%v subject=%q", m.To, m.Subject)
		}
	}()
}
