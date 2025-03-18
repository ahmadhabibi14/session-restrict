package mailer

import (
	"os"
	"session-restrict/src/lib/logger"
	"strconv"
	"strings"

	"github.com/wneessen/go-mail"
)

type Mailhog struct {
	client *mail.Client
}

func NewMailhog() (*Mailhog, error) {
	port, err := strconv.Atoi(os.Getenv(`MAILHOG_PORT`))
	if err != nil {
		port = 456
	}

	mailhog := &Mailhog{}

	client, err := mail.NewClient(
		os.Getenv(`MAILHOG_HOST`),
		mail.WithPort(port),
		mail.WithTLSPolicy(mail.NoTLS),
	)
	if err != nil {
		logger.Log.Error(err)
		return mailhog, err
	}

	mailhog.client = client

	return mailhog, nil
}

func (m *Mailhog) SendMailText(to []string, cc []string, subject, message string) error {
	msg := mail.NewMsg()
	err := msg.FromFormat(
		os.Getenv(`SMTP_SENDER_NAME`),
		os.Getenv(`SMTP_USERNAME`),
	)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	err = msg.To(to...)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	err = msg.Cc(cc...)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, message)
	msg.SetImportance(mail.ImportanceHigh)

	err = m.client.DialAndSend(msg)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	logger.Log.Info(`email sent to: ` + strings.Join(to, ","))

	return nil
}

func (m *Mailhog) SendMailHTML(to []string, cc []string, subject, htmlString string) error {
	msg := mail.NewMsg()
	err := msg.FromFormat(
		os.Getenv(`SMTP_SENDER_NAME`),
		os.Getenv(`SMTP_USERNAME`),
	)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}
	err = msg.To(to...)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}
	err = msg.Cc(cc...)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, htmlString)
	msg.SetImportance(mail.ImportanceHigh)

	err = m.client.DialAndSend(msg)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}
	logger.Log.Info(`email sent to: ` + strings.Join(to, ","))
	return nil
}
