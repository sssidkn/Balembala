package repository

import (
	"context"
	"database/sql"
	"errors"
	"report/pkg/api/report"
	"report/pkg/postgres"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestCreateTemplate(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		input       *report.Template
		userId      int64
		expected    int64
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE user_id = \$1 AND title = \$2\)`).
					WithArgs(int64(1), "Test Template").
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
				mock.ExpectQuery(`INSERT INTO templates \(title, message, user_id\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
					WithArgs("Test Template", "Test message", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			},
			input: &report.Template{
				Title:   "Test Template",
				Message: "Test message",
			},
			userId:      1,
			expected:    1,
			expectedErr: nil,
		},
		{
			name: "Template already exists",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE user_id = \$1 AND title = \$2\)`).
					WithArgs(int64(1), "Existing Template").
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
			},
			input: &report.Template{
				Title:   "Existing Template",
				Message: "Test message",
			},
			userId:      1,
			expected:    0,
			expectedErr: ErrTemplateAlreadyExists,
		},
		{
			name: "Database error on existence check",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE user_id = \$1 AND title = \$2\)`).
					WithArgs(int64(1), "Test Template").
					WillReturnError(errors.New("db error"))
			},
			input: &report.Template{
				Title:   "Test Template",
				Message: "Test message",
			},
			userId:      1,
			expected:    0,
			expectedErr: errors.New("failed to check existence: db error"),
		},
		{
			name: "Database error on insert",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE user_id = \$1 AND title = \$2\)`).
					WithArgs(int64(1), "Test Template").
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
				mock.ExpectQuery(`INSERT INTO templates \(title, message, user_id\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
					WithArgs("Test Template", "Test message", int64(1)).
					WillReturnError(errors.New("db error"))
			},
			input: &report.Template{
				Title:   "Test Template",
				Message: "Test message",
			},
			userId:      1,
			expected:    0,
			expectedErr: errors.New("failed to create template: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.CreateTemplate(context.Background(), tt.input, tt.userId)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				return
			}
			if got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		templateId  int64
		userId      int64
		expected    *report.Template
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id", "title", "message"}).
					AddRow(int64(1), "Test Template", "Test message")
				mock.ExpectQuery(`SELECT id, title, message FROM templates WHERE id = \$1 AND user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(rows)
			},
			templateId: 1,
			userId:     1,
			expected: &report.Template{
				Id:      1,
				Title:   "Test Template",
				Message: "Test message",
			},
			expectedErr: nil,
		},
		{
			name: "Not found",
			mock: func() {
				mock.ExpectQuery(`SELECT id, title, message FROM templates WHERE id = \$1 AND user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			templateId:  1,
			userId:      1,
			expected:    nil,
			expectedErr: ErrNotFound,
		},
		{
			name: "Database error",
			mock: func() {
				mock.ExpectQuery(`SELECT id, title, message FROM templates WHERE id = \$1 AND user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnError(errors.New("db error"))
			},
			templateId:  1,
			userId:      1,
			expected:    nil,
			expectedErr: errors.New("failed to get template: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetTemplate(context.Background(), tt.templateId, tt.userId)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.expectedErr != nil && (err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				return
			}
			if got != nil && tt.expected != nil {
				if got.Id != tt.expected.Id || got.Title != tt.expected.Title || got.Message != tt.expected.Message {
					t.Errorf("expected: %v, got: %v", tt.expected, got)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTemplate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		input       *report.Template
		userId      int64
		expected    *report.Template
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectQuery(`UPDATE templates SET title = \$1, message = \$2 WHERE id = \$3 RETURNING id, title, message`).
					WithArgs("Updated", "New message", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "title", "message"}).
						AddRow(1, "Updated", "New message"))
			},
			input: &report.Template{
				Id:      1,
				Title:   "Updated",
				Message: "New message",
			},
			userId: 1,
			expected: &report.Template{
				Id:      1,
				Title:   "Updated",
				Message: "New message",
			},
			expectedErr: nil,
		},
		{
			name: "Template not found",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
			},
			input: &report.Template{
				Id:      1,
				Title:   "Updated",
				Message: "New message",
			},
			userId:      1,
			expected:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.UpdateTemplate(context.Background(), tt.input, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != nil && tt.expected != nil {
				if result.Id != tt.expected.Id || result.Title != tt.expected.Title || result.Message != tt.expected.Message {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func TestDeleteTemplate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		templateId  int64
		userId      int64
		expected    bool
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectExec(`DELETE FROM templates WHERE id = \$1`).
					WithArgs(int64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			templateId:  1,
			userId:      1,
			expected:    true,
			expectedErr: nil,
		},
		{
			name: "Template not found",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
			},
			templateId:  1,
			userId:      1,
			expected:    false,
			expectedErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.DeleteTemplate(context.Background(), tt.templateId, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		userId      int64
		expected    []*report.Template
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id", "title", "message"}).
					AddRow(int64(1), "Template 1", "Message 1").
					AddRow(int64(2), "Template 2", "Message 2")
				mock.ExpectQuery(`SELECT id, title, message FROM templates WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			userId: 1,
			expected: []*report.Template{
				{Id: 1, Title: "Template 1", Message: "Message 1"},
				{Id: 2, Title: "Template 2", Message: "Message 2"},
			},
			expectedErr: nil,
		},
		{
			name: "No templates",
			mock: func() {
				mock.ExpectQuery(`SELECT id, title, message FROM templates WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			userId:      1,
			expected:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.ListTemplates(context.Background(), tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d templates, got %d", len(tt.expected), len(result))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCreateContact(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		input       *report.Contact
		userId      int64
		expected    int64
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contacts WHERE email = \$1 AND user_id = \$2\)`).
					WithArgs("test@example.com", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
				mock.ExpectQuery(`INSERT INTO contacts \(name, email, user_id\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
					WithArgs("Test", "test@example.com", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			},
			input: &report.Contact{
				Name:  "Test",
				Email: "test@example.com",
			},
			userId:      1,
			expected:    1,
			expectedErr: nil,
		},
		{
			name: "Contact already exists",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contacts WHERE email = \$1 AND user_id = \$2\)`).
					WithArgs("test@example.com", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
			},
			input: &report.Contact{
				Name:  "Test",
				Email: "test@example.com",
			},
			userId:      1,
			expected:    0,
			expectedErr: ErrContactAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.CreateContact(context.Background(), tt.input, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetContact(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		contactId   int64
		userId      int64
		expected    *report.Contact
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id", "email", "name"}).
					AddRow(1, "test@example.com", "Test")
				mock.ExpectQuery(`SELECT id, email, name FROM contacts WHERE id = \$1 AND user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(rows)
			},
			contactId: 1,
			userId:    1,
			expected: &report.Contact{
				Id:    1,
				Email: "test@example.com",
				Name:  "Test",
			},
			expectedErr: nil,
		},
		{
			name: "Contact not found",
			mock: func() {
				mock.ExpectQuery(`SELECT id, email, name FROM contacts WHERE id = \$1 AND user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			contactId:   1,
			userId:      1,
			expected:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.GetContact(context.Background(), tt.contactId, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != nil && tt.expected != nil {
				if result.Id != tt.expected.Id || result.Email != tt.expected.Email || result.Name != tt.expected.Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateContact(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		input       *report.Contact
		userId      int64
		expected    *report.Contact
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contacts WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectQuery(`UPDATE contacts SET email = \$1, name = \$2 WHERE id = \$3 RETURNING email, name, id`).
					WithArgs("new@example.com", "Updated", int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "name", "id"}).
						AddRow("new@example.com", "Updated", int64(1)))
			},
			input: &report.Contact{
				Id:    1,
				Email: "new@example.com",
				Name:  "Updated",
			},
			userId: 1,
			expected: &report.Contact{
				Id:    1,
				Email: "new@example.com",
				Name:  "Updated",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.UpdateContact(context.Background(), tt.input, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != nil && tt.expected != nil {
				if result.Id != tt.expected.Id || result.Email != tt.expected.Email || result.Name != tt.expected.Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteContact(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		contactId   int64
		userId      int64
		expected    bool
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contacts WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectExec(`DELETE FROM contacts WHERE id = \$1`).
					WithArgs(int64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			contactId:   1,
			userId:      1,
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.DeleteContact(context.Background(), tt.contactId, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetContactsByTemplate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		templateId  int64
		userId      int64
		expected    []*report.Contact
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id", "email", "name"}).
					AddRow(int64(1), "test1@example.com", "Test 1").
					AddRow(int64(2), "test2@example.com", "Test 2")
				mock.ExpectQuery(`SELECT c.id, c.email, c.name FROM contacts c JOIN contacts_templates ct ON c.id = ct.contact_id WHERE ct.template_id = \$1 AND c.user_id = \$2`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(rows)
			},
			templateId: 1,
			userId:     1,
			expected: []*report.Contact{
				{Id: 1, Email: "test1@example.com", Name: "Test 1"},
				{Id: 2, Email: "test2@example.com", Name: "Test 2"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.GetContactsByTemplate(context.Background(), tt.templateId, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d contacts, got %d", len(tt.expected), len(result))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddContactsToTemplate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		templateId  int64
		contactIds  []int64
		userId      int64
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM templates WHERE id = \$1 AND user_id = \$2\)`).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM contacts WHERE id = ANY\(\$1\) AND user_id = \$2`).
					WithArgs([]int64{1, 2}, int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(2)))
				mock.ExpectExec(`INSERT INTO contacts_templates \(template_id, contact_id\) VALUES \(\$1, \$2\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(`INSERT INTO contacts_templates \(template_id, contact_id\) VALUES \(\$1, \$2\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(2)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
				mock.ExpectRollback()
			},
			templateId:  1,
			contactIds:  []int64{1, 2},
			userId:      1,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.AddContactsToTemplate(context.Background(), tt.templateId, tt.contactIds, tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetContacts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &ReportRepository{db: &postgres.DB{Db: mock}}

	tests := []struct {
		name        string
		mock        func()
		userId      int64
		expected    []*report.Contact
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id", "name", "email"}).
					AddRow(int64(1), "Test 1", "test1@example.com").
					AddRow(int64(2), "Test 2", "test2@example.com")
				mock.ExpectQuery(`SELECT id, name, email FROM contacts WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			userId: 1,
			expected: []*report.Contact{
				{Id: 1, Name: "Test 1", Email: "test1@example.com"},
				{Id: 2, Name: "Test 2", Email: "test2@example.com"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.GetContacts(context.Background(), tt.userId)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d contacts, got %d", len(tt.expected), len(result))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
