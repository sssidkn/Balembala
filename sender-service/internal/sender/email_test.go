package sender

import (
	"crypto/tls"
	"errors"
	"io"
	"net/smtp"
	"sender-service/internal/dto"
	"testing"
)

// MockSMTPClient реализует SMTPClient интерфейс для тестирования
type MockSMTPClient struct {
	StartTLSFunc    func(config *tls.Config) error
	AuthFunc        func(auth smtp.Auth) error
	MailFunc        func(from string) error
	RcptFunc        func(to string) error
	DataFunc        func() (io.WriteCloser, error)
	QuitFunc        func() error
	CloseFunc       func() error
	RcptFailFor     map[string]bool // Адреса, для которых Rcpt должен вернуть ошибку
	DataShouldFail  bool
	WriteShouldFail bool
	CloseShouldFail bool
}

func (m *MockSMTPClient) StartTLS(config *tls.Config) error {
	if m.StartTLSFunc != nil {
		return m.StartTLSFunc(config)
	}
	return nil
}

func (m *MockSMTPClient) Auth(auth smtp.Auth) error {
	if m.AuthFunc != nil {
		return m.AuthFunc(auth)
	}
	return nil
}

func (m *MockSMTPClient) Mail(from string) error {
	if m.MailFunc != nil {
		return m.MailFunc(from)
	}
	return nil
}

func (m *MockSMTPClient) Rcpt(to string) error {
	if m.RcptFunc != nil {
		return m.RcptFunc(to)
	}
	if m.RcptFailFor != nil && m.RcptFailFor[to] {
		return errors.New("rcpt failed for " + to)
	}
	return nil
}

func (m *MockSMTPClient) Data() (io.WriteCloser, error) {
	if m.DataFunc != nil {
		return m.DataFunc()
	}
	if m.DataShouldFail {
		return nil, errors.New("DATA command failed")
	}
	return &mockWriteCloser{
		WriteShouldFail: m.WriteShouldFail,
		CloseShouldFail: m.CloseShouldFail,
	}, nil
}

func (m *MockSMTPClient) Quit() error {
	if m.QuitFunc != nil {
		return m.QuitFunc()
	}
	return nil
}

func (m *MockSMTPClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type mockWriteCloser struct {
	WriteShouldFail bool
	CloseShouldFail bool
	WrittenData     []byte
}

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	if m.WriteShouldFail {
		return 0, errors.New("write failed")
	}
	m.WrittenData = append(m.WrittenData, p...)
	return len(p), nil
}

func (m *mockWriteCloser) Close() error {
	if m.CloseShouldFail {
		return errors.New("close failed")
	}
	return nil
}

// MockDialer реализует SMTPDialer для тестирования
func MockDialer(client SMTPClient) SMTPDialer {
	return func(addr string) (SMTPClient, error) {
		return client, nil
	}
}

func TestEmail_Send(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockSMTPClient
		message       dto.Message
		wantErr       bool
		wantRetryList []string
	}{
		{
			name: "successful send to all recipients",
			client: &MockSMTPClient{
				RcptFailFor: make(map[string]bool),
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com", "user2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       false,
			wantRetryList: []string{},
		},
		{
			name: "failed rcpt for one recipient",
			client: &MockSMTPClient{
				RcptFailFor: map[string]bool{"user2@example.com": true},
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com", "user2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       false,
			wantRetryList: []string{"user2@example.com"},
		},
		{
			name: "failed DATA command",
			client: &MockSMTPClient{
				DataShouldFail: true,
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com", "user2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       true,
			wantRetryList: []string{"user1@example.com", "user2@example.com"},
		},
		{
			name: "failed Write",
			client: &MockSMTPClient{
				WriteShouldFail: true,
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com", "user2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       true,
			wantRetryList: []string{"user1@example.com", "user2@example.com"},
		},
		{
			name: "failed Close",
			client: &MockSMTPClient{
				CloseShouldFail: true,
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com", "user2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       true,
			wantRetryList: []string{"user1@example.com", "user2@example.com"},
		},
		{
			name: "failed Mail command",
			client: &MockSMTPClient{
				MailFunc: func(from string) error {
					return errors.New("mail failed")
				},
			},
			message: dto.Message{
				ToList:  []string{"user1@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:       true,
			wantRetryList: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{
				Username:  "sender@example.com",
				Password:  "password",
				Port:      587,
				Host:      "smtp.example.com",
				tlsConfig: &tls.Config{},
				dialer:    MockDialer(tt.client),
			}

			err, retryMessage := e.Send(tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("Email.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(retryMessage.ToList) != len(tt.wantRetryList) {
				t.Errorf("Email.Send() retry list length = %d, want %d", len(retryMessage.ToList), len(tt.wantRetryList))
			}

			if tt.wantErr && len(retryMessage.ToList) != len(tt.message.ToList) {
				t.Errorf("Email.Send() on error should return all recipients for retry, got %d, want %d",
					len(retryMessage.ToList), len(tt.message.ToList))
			}

			// Проверяем, что все ожидаемые адреса есть в списке повторной отправки
			for _, wantAddr := range tt.wantRetryList {
				found := false
				for _, gotAddr := range retryMessage.ToList {
					if gotAddr == wantAddr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Email.Send() missing address in retry list: %s", wantAddr)
				}
			}
		})
	}
}

func TestEmail_Send_ConnectionError(t *testing.T) {
	e := &Email{
		Username:  "sender@example.com",
		Password:  "password",
		Port:      587,
		Host:      "smtp.example.com",
		tlsConfig: &tls.Config{},
		dialer: func(addr string) (SMTPClient, error) {
			return nil, errors.New("connection failed")
		},
	}

	message := dto.Message{
		ToList:  []string{"user@example.com"},
		Subject: "Test",
		Body:    "Test",
	}

	err, retryMessage := e.Send(message)
	if err == nil {
		t.Error("Expected error for failed connection, got nil")
	}
	if len(retryMessage.ToList) == 0 {
		t.Error("Expected full retry list for connection error")
	}
}

func TestEmail_Send_AuthError(t *testing.T) {
	client := &MockSMTPClient{
		AuthFunc: func(auth smtp.Auth) error {
			return errors.New("auth failed")
		},
	}

	e := &Email{
		Username:  "sender@example.com",
		Password:  "password",
		Port:      587,
		Host:      "smtp.example.com",
		tlsConfig: &tls.Config{},
		dialer:    MockDialer(client),
	}

	message := dto.Message{
		ToList:  []string{"user@example.com"},
		Subject: "Test",
		Body:    "Test",
	}

	err, retryMessage := e.Send(message)
	if err == nil {
		t.Error("Expected error for failed auth, got nil")
	}
	if len(retryMessage.ToList) == 0 {
		t.Error("Expected full retry list for auth error")
	}
}
