package util

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

func ValidateEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

type EmailSender struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort int
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

func (s *EmailSender) Send(msg EmailMessage) error {
	auth := smtp.PlainAuth("", s.From, s.Password, s.SMTPHost)

	message := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		msg.To, msg.Subject, msg.Body,
	))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort),
		auth,
		s.From,
		[]string{msg.To},
		message,
	)

	return err
}
