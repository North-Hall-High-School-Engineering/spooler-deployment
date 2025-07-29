package util

import (
	"fmt"
	"net/smtp"
	"regexp"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

type EmailSender struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
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
		s.SMTPHost+":"+s.SMTPPort,
		auth,
		s.From,
		[]string{msg.To},
		message,
	)

	return err
}
