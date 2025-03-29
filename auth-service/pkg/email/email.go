package email

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	smtpHost string
	smtpPort string
	auth     smtp.Auth
	from     string
}

func NewMailService(smtpHost, smtpPort, username, password, from string) *MailService {
	auth := smtp.PlainAuth("", username, password, smtpHost)
	return &MailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		auth:     auth,
		from:     from,
	}
}
func (ms *MailService) SendWelcomeEmail(toEmail, userName string) error {
	subject := "Welcome to Our Service"
	body := fmt.Sprintf(
		"Hello %s,\n\nThank you for registering with us. We're excited to have you on board!\n\nBest regards,\nThe Team",
		userName,
	)
	message := fmt.Sprintf("From: %s\r\n", ms.from) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + body
	addr := fmt.Sprintf("%s:%s", ms.smtpHost, ms.smtpPort)
	if err := smtp.SendMail(addr, ms.auth, ms.from, []string{toEmail}, []byte(message)); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}
