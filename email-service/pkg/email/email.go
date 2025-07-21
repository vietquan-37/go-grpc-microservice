package email

import (
	"bytes"
	"fmt"
	"github.com/vietquan-37/email-service/pkg/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
)

type IEmailService interface {
	SendVerificationEmail(to, fullName, token string, emailType Type) error
	SendOrderConfirmationEmail(to, fullName string, orderData OrderData, emailType Type) error
}

type MailService struct {
	config    *config.Config
	dialer    *gomail.Dialer
	templates *template.Template
}
type BaseEmailData struct {
	FullName     string
	CompanyName  string
	SupportEmail string
	Year         int
}
type VerificationEmail struct {
	BaseEmailData
	ActiveLink string
}
type OrderData struct {
	OrderID     int32
	OrderDate   string
	Status      string
	TotalAmount float64
	Items       []OrderItem
}

type OrderItem struct {
	Name     string
	Quantity int64
	Price    float64
}

type OrderConfirmationEmailData struct {
	BaseEmailData
	Order OrderData
}

func NewEmailService(config *config.Config) (IEmailService, error) {
	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rootDir := filepath.Dir(wd)
	// Build absolute path to templates
	templatePath := filepath.Join(rootDir, "pkg", "templates", "*.html")
	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, err
	}

	return &MailService{
		config:    config,
		dialer:    dialer,
		templates: templates,
	}, nil
}
func (e *MailService) createBaseData(fullName string) BaseEmailData {
	return BaseEmailData{
		FullName:     fullName,
		CompanyName:  "Our Service",
		SupportEmail: e.config.SMTPUsername,
		Year:         2025,
	}
}
func (e *MailService) SendVerificationEmail(to, fullName, token string, emailType Type) error {

	activationLink := fmt.Sprintf("%s/activate?token=%s", e.config.URL, token)
	t, ok := Templates[emailType]
	if !ok {
		return fmt.Errorf("unknown email type: %s", emailType)
	}
	data := VerificationEmail{
		BaseEmailData: e.createBaseData(fullName),
		ActiveLink:    activationLink,
	}

	return e.sendTemplateEmail(to, t.Subject, t.TemplateName, data)
}
func (e *MailService) SendOrderConfirmationEmail(to, fullName string, orderData OrderData, emailType Type) error {
	t, ok := Templates[emailType]
	if !ok {
		return fmt.Errorf("unknown email type: %s", emailType)
	}
	data := OrderConfirmationEmailData{
		BaseEmailData: e.createBaseData(fullName),
		Order:         orderData,
	}
	return e.sendTemplateEmail(to, t.Subject, t.TemplateName, data)
}
func (e *MailService) sendTemplateEmail(to, subject, templateName string, data interface{}) error {
	var body bytes.Buffer

	err := e.templates.ExecuteTemplate(&body, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.config.SMTPUsername)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body.String())

	return e.dialer.DialAndSend(msg)
}
