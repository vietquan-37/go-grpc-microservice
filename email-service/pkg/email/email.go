package email

import (
	"fmt"

	"github.com/vietquan-37/email-service/pkg/config"
	"gopkg.in/gomail.v2"
)

type IEmailService interface {
	SendVerificationEmail(to, fullName, token string) error
}

type MailService struct {
	config *config.Config
	dialer *gomail.Dialer
}

func NewEmailService(config *config.Config) IEmailService {
	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	return &MailService{
		config: config,
		dialer: dialer,
	}
}

func (e *MailService) SendVerificationEmail(to, fullName, token string) error {
	subject := "Welcome to Our Service!"

	activationLink := fmt.Sprintf("%s/activate?token=%s", e.config.URL, token)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome %s!</h1>
			<p>Thank you for registering with our service.</p>
			<p>Your account has been created successfully.</p>
			<p><a href="%s">Click here to activate your account</a></p>
		</body>
		</html>
	`, fullName, activationLink)

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.config.SMTPUsername)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return e.dialer.DialAndSend(msg)
}
