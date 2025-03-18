package mailer

import (
	"errors"
	"os"
	"session-restrict/configs"
)

type FuncSendMailText func(to []string, cc []string, subject, message string) error
type FuncSendMailHTML func(to []string, cc []string, subject, htmlString string) error

var ErrSendMail = errors.New(`failed to send email`)

type Mailer struct {
	FuncSendMailText FuncSendMailText
	FuncSendMailHTML FuncSendMailHTML
}

func NewMailer() (*Mailer, error) {
	var mailer = &Mailer{}
	var projectEnv = os.Getenv(`PROJECT_ENV`)

	switch projectEnv {
	case configs.EnvDev:
		ml, err := NewMailhog()
		if err != nil {
			return mailer, err
		}

		mailer.FuncSendMailText = ml.SendMailText
		mailer.FuncSendMailHTML = ml.SendMailHTML
	default:
		// Replace this with actual SMTP Server
		ml, err := NewMailhog()
		if err != nil {
			return mailer, err
		}

		mailer.FuncSendMailText = ml.SendMailText
		mailer.FuncSendMailHTML = ml.SendMailHTML
	}

	return mailer, nil
}

func (m *Mailer) SendMailText(to []string, cc []string, subject, message string) error {
	return m.FuncSendMailText(to, cc, subject, message)
}

func (m *Mailer) SendMailHTML(to []string, cc []string, subject, htmlString string) error {
	return m.FuncSendMailHTML(to, cc, subject, htmlString)
}
