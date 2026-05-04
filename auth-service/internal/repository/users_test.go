package repository

import (
	"auth/pkg/db/postgres"
	"context"
	"errors"
	"testing"

	"auth/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_Register(t *testing.T) {
	tests := []struct {
		name        string
		mock        func(mock pgxmock.PgxPoolIface)
		input       *models.RegisterRequest
		expected    *models.User
		expectedErr error
	}{
		{
			name: "Success",
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("test@example.com", "testuser", "hashedpassword").
					WillReturnRows(rows)
			},
			input: &models.RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "hashedpassword",
			},
			expected: &models.User{
				UserID:   1,
				Email:    "test@example.com",
				Username: "testuser",
			},
			expectedErr: nil,
		},
		{
			name: "Duplicate email",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("exists@example.com", "testuser", "hashedpassword").
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			input: &models.RegisterRequest{
				Email:    "exists@example.com",
				Username: "testuser",
				Password: "hashedpassword",
			},
			expected:    nil,
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name: "Other database error",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("error@example.com", "testuser", "hashedpassword").
					WillReturnError(errors.New("database error"))
			},
			input: &models.RegisterRequest{
				Email:    "error@example.com",
				Username: "testuser",
				Password: "hashedpassword",
			},
			expected:    nil,
			expectedErr: errors.New("repository.Register: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			repo := &UserRepository{db: &postgres.DB{Db: mock}}
			tt.mock(mock)

			result, err := repo.Register(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_Login(t *testing.T) {
	validPassword := "correctpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		mock        func(mock pgxmock.PgxPoolIface)
		input       *models.LoginRequest
		expected    *models.User
		expectedErr error
	}{
		{
			name: "Success",
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"username", "password_hash", "id"}).
					AddRow("testuser", string(hashedPassword), 1)
				mock.ExpectQuery("SELECT username, password_hash, id FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			input: &models.LoginRequest{
				Email:    "test@example.com",
				Password: validPassword,
			},
			expected: &models.User{
				UserID:   1,
				Email:    "test@example.com",
				Username: "testuser",
			},
			expectedErr: nil,
		},
		{
			name: "User not found",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT username, password_hash, id FROM users WHERE email = \\$1").
					WithArgs("notfound@example.com").
					WillReturnError(errors.New("no rows in result set"))
			},
			input: &models.LoginRequest{
				Email:    "notfound@example.com",
				Password: "anypassword",
			},
			expected:    nil,
			expectedErr: ErrUserNotFound,
		},
		{
			name: "Invalid password",
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"username", "password_hash", "id"}).
					AddRow("testuser", string(hashedPassword), 1)
				mock.ExpectQuery("SELECT username, password_hash, id FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			input: &models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expected:    nil,
			expectedErr: ErrInvalidPassword,
		},
		{
			name: "Database error",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT username, password_hash, id FROM users WHERE email = \\$1").
					WithArgs("error@example.com").
					WillReturnError(errors.New("database error"))
			},
			input: &models.LoginRequest{
				Email:    "error@example.com",
				Password: "anypassword",
			},
			expected:    nil,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			repo := &UserRepository{db: &postgres.DB{Db: mock}}
			tt.mock(mock)

			result, err := repo.Login(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestIsDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{
			name:     "Is duplicate key",
			input:    &pgconn.PgError{Code: "23505"},
			expected: true,
		},
		{
			name:     "Not a pg error",
			input:    errors.New("regular error"),
			expected: false,
		},
		{
			name:     "Different pg error code",
			input:    &pgconn.PgError{Code: "12345"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDuplicateKeyError(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
