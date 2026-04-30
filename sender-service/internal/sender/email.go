package sender

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"sender-service/internal/config"
	"sender-service/internal/dto"
	"strings"
)

type SMTPClient interface {
	StartTLS(config *tls.Config) error
	Auth(auth smtp.Auth) error
	Mail(from string) error
	Rcpt(to string) error
	Data() (io.WriteCloser, error)
	Quit() error
	Close() error
}

type SMTPDialer func(addr string) (SMTPClient, error)

type Sender interface {
	Send(dto.Message) error
}

type Email struct {
	Username   string
	Password   string
	Port       int
	Host       string
	tlsConfig  *tls.Config
	caCertPath string
	dialer     SMTPDialer
}

func New(cfg config.Config) *Email {
	return NewWithDialer(cfg, defaultDialer)
}

func NewWithDialer(cfg config.Config, dialer SMTPDialer) *Email {
	email := &Email{
		Port:       587,
		Host:       "smtp.mail.ru",
		Username:   cfg.Username,
		Password:   cfg.Password,
		caCertPath: cfg.CaCertPath,
		dialer:     dialer,
	}
	email.initTLSConfig()
	return email
}

func (e *Email) initTLSConfig() {
	e.tlsConfig = &tls.Config{
		ServerName:         e.Host,
		InsecureSkipVerify: true,
	}
}

func (e *Email) Send(message dto.Message) (error, dto.Message) {
	var retryMessage dto.Message
	client, err := e.dialer(fmt.Sprintf("%s:%d", e.Host, e.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err), message
	}
	defer client.Close()

	if err = client.StartTLS(e.tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err), message
	}

	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err), message
	}

	if err = client.Mail(e.Username); err != nil {
		return fmt.Errorf("mail from failed: %w", err), message
	}

	for _, to := range message.ToList {
		if err = client.Rcpt(to); err != nil {
			retryMessage.ToList = append(retryMessage.ToList, to)
		}
	}
	retryMessage.Subject = message.Subject
	retryMessage.Body = message.Body

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err), message
	}

	msg := []byte(
		"From: " + e.Username + "\r\n" +
			"To: " + strings.Join(message.ToList, ",") + "\r\n" +
			"Subject: " + message.Subject + "\r\n" +
			"\r\n" + message.Body,
	)

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err), message
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err), message
	}
	client.Quit()
	return nil, retryMessage
}

func defaultDialer(addr string) (SMTPClient, error) {
	return smtp.Dial(addr)
}
