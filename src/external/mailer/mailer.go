package mailer

import (
	"errors"
	"os"
	"session-restrict/src/lib/logger"
	"strconv"
	"strings"

	"github.com/wneessen/go-mail"
)

type MailerService struct {
	client *mail.Client
}

func NewMailerService() *MailerService {
	port, err := strconv.Atoi(os.Getenv(`SMTP_PORT`))
	if err != nil {
		port = 456
	}

	client, err := mail.NewClient(
		os.Getenv(`SMTP_HOST`),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithSSL(),
		mail.WithPort(port),
		mail.WithUsername(os.Getenv(`SMTP_USERNAME`)),
		mail.WithPassword(os.Getenv(`SMTP_PASSWORD`)),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		logger.Log.Fatal(err, `failed to initialize mail client`)
	}

	logger.Log.Info(`Initialized mail client`)

	return &MailerService{
		client: client,
	}
}

var (
	ErrSendMail = errors.New(`failed to send email`)
)

func (ms *MailerService) SendMailText(to []string, cc []string, subject, message string) error {
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

	err = ms.client.DialAndSend(msg)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	logger.Log.Info(`Email sent to: ` + strings.Join(to, ","))

	return nil
}

func (ms *MailerService) SendMailHTML(to []string, cc []string, subject, htmlString string) error {
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

	err = ms.client.DialAndSend(msg)
	if err != nil {
		logger.Log.Error(err, ErrSendMail.Error())
		return ErrSendMail
	}

	logger.Log.Info(`Email sent to: ` + strings.Join(to, ","))

	return nil
}
