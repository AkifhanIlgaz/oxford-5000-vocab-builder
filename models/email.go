package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@vocab-builder.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	service := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.UserName, config.Password),
	}

	return &service
}

func (service *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	service.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.PlainText != "" && email.HTML != "":
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		msg.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	if err := service.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil

}

func (service *EmailService) ForgotPassword(to, resetUrl string) error {
	email := Email{
		To:        to,
		Subject:   "Reset your password",
		PlainText: fmt.Sprintf("To reset your password, please visit the following link: %s", resetUrl),
		HTML:      fmt.Sprintf(`<p>To reset your password, please visit the following link: <a href="%s">%s</a></p>`, resetUrl, resetUrl),
	}

	if err := service.Send(email); err != nil {
		return fmt.Errorf("forgot passwor email: %w", err)
	}

	return nil
}

func (service *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string

	switch {
	case email.From != "":
		from = email.From
	case service.DefaultSender != "":
		from = service.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}
