package workout

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailService interface {
	SendWorkoutEmail(to string, subject string, htmlBody string) error
}

type emailService struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
	fromName  string
}

func NewEmailService() EmailService {
	return &emailService{
		host:      os.Getenv("SMTP_HOST"),
		port:      os.Getenv("SMTP_PORT"),
		username:  os.Getenv("SMTP_USERNAME"),
		password:  os.Getenv("SMTP_PASSWORD"),
		fromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		fromName:  os.Getenv("SMTP_FROM_NAME"),
	}
}

func (e *emailService) SendWorkoutEmail(to string, subject string, htmlBody string) error {
	if e.host == "" || e.port == "" || e.username == "" || e.password == "" || e.fromEmail == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	from := e.fromEmail
	if strings.TrimSpace(e.fromName) != "" {
		from = fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail)
	}

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
		"",
		htmlBody,
	}, "\r\n")

	address := e.host + ":" + e.port
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	return smtp.SendMail(
		address,
		auth,
		e.fromEmail,
		[]string{to},
		[]byte(message),
	)
}
