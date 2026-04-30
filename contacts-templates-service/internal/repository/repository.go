package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"report/pkg/api/report"
	"report/pkg/postgres"
)

//go:generate mockgen -source=irepository.go -destination=repository_mock.go -package=repository

var (
	ErrTemplateAlreadyExists = errors.New("template already exists")
	ErrContactAlreadyExists  = errors.New("contact already exists")
	ErrNotFound              = errors.New("not found")
)

type Repository interface {
	CreateTemplate(ctx context.Context, t *report.Template, userId int64) (int64, error)
	GetTemplate(ctx context.Context, templateId int64, userId int64) (*report.Template, error)
	UpdateTemplate(ctx context.Context, t *report.Template, userId int64) (*report.Template, error)
	DeleteTemplate(ctx context.Context, templateId int64, userId int64) (bool, error)
	ListTemplates(ctx context.Context, userId int64) ([]*report.Template, error)
	CreateContact(ctx context.Context, c *report.Contact, userId int64) (int64, error)
	GetContact(ctx context.Context, contactId int64, userId int64) (*report.Contact, error)
	UpdateContact(ctx context.Context, c *report.Contact, userId int64) (*report.Contact, error)
	DeleteContact(ctx context.Context, contactId int64, userId int64) (bool, error)
	GetContactsByTemplate(ctx context.Context, templateId int64, userId int64) ([]*report.Contact, error)
	AddContactsToTemplate(ctx context.Context, templateId int64, contactIds []int64, userId int64) error
	GetContacts(ctx context.Context, userId int64) ([]*report.Contact, error)
}

type ReportRepository struct {
	db *postgres.DB
}

func NewRepository(db *postgres.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateTemplate(ctx context.Context, t *report.Template, userId int64) (int64, error) {
	conn := r.db.Db
	var id int64
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM templates WHERE user_id = $1 AND title = $2)`,
		userId, t.Title).Scan(&exist)
	if err != nil {
		return 0, fmt.Errorf("failed to check existence: %w", err)
	}
	if exist {
		return id, ErrTemplateAlreadyExists
	}
	err = conn.QueryRow(ctx, `INSERT INTO templates (title, message, user_id)
		VALUES ($1, $2, $3) RETURNING id`, t.Title, t.Message, userId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create template: %w", err)
	}
	return id, nil
}

func (r *ReportRepository) GetTemplate(ctx context.Context, templateId int64, userId int64) (*report.Template, error) {
	conn := r.db.Db
	var template report.Template
	err := conn.QueryRow(ctx, `SELECT id, title, message FROM templates WHERE id = $1 AND user_id = $2`,
		templateId, userId).Scan(&template.Id, &template.Title, &template.Message)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &template, nil
}

func (r *ReportRepository) UpdateTemplate(ctx context.Context, t *report.Template, userId int64) (*report.Template, error) {
	conn := r.db.Db
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM templates WHERE id = $1 AND user_id = $2)`, t.Id, userId).Scan(&exist)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exist {
		return nil, ErrNotFound
	}
	err = conn.QueryRow(ctx, `UPDATE templates SET title = $1, message = $2 
                 WHERE id = $3 RETURNING id, title, message`, t.Title, t.Message,
		t.Id).Scan(&t.Id, &t.Title, &t.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}
	return t, nil
}

func (r *ReportRepository) DeleteTemplate(ctx context.Context, templateId int64, userId int64) (bool, error) {
	conn := r.db.Db
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM templates WHERE id = $1 AND user_id = $2)`,
		templateId, userId).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exist {
		return false, ErrNotFound
	}
	_, err = conn.Exec(ctx, `DELETE FROM templates WHERE id = $1`, templateId)
	if err != nil {
		return false, fmt.Errorf("failed to delete template: %w", err)
	}
	return true, nil
}

func (r *ReportRepository) ListTemplates(ctx context.Context, userId int64) ([]*report.Template, error) {
	conn := r.db.Db
	var templates []*report.Template
	rows, err := conn.Query(ctx, `SELECT id, title, message FROM templates WHERE user_id = $1`, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t report.Template
		err = rows.Scan(&t.Id, &t.Title, &t.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (r *ReportRepository) CreateContact(ctx context.Context, c *report.Contact, userId int64) (int64, error) {
	conn := r.db.Db
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contacts WHERE email = $1 AND user_id = $2)`,
		c.Email, userId).Scan(&exist)
	if err != nil {
		return 0, fmt.Errorf("failed to check existence: %w", err)
	}
	if exist {
		return 0, ErrContactAlreadyExists
	}
	err = conn.QueryRow(ctx, `INSERT INTO contacts (name, email, user_id) VALUES ($1, $2, $3) RETURNING id`,
		c.Name, c.Email, userId).Scan(&c.Id)
	if err != nil {
		return 0, fmt.Errorf("failed to create contact: %w", err)
	}
	return c.Id, nil
}

func (r *ReportRepository) GetContact(ctx context.Context, contactId int64, userId int64) (*report.Contact, error) {
	conn := r.db.Db
	var contact report.Contact
	err := conn.QueryRow(ctx, `SELECT id, email, name FROM contacts WHERE id = $1 AND user_id = $2`, contactId,
		userId).Scan(&contact.Id, &contact.Email, &contact.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}
	return &contact, nil
}

func (r *ReportRepository) UpdateContact(ctx context.Context, c *report.Contact, userId int64) (*report.Contact, error) {
	conn := r.db.Db
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`,
		c.Id, userId).Scan(&exist)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exist {
		return nil, ErrNotFound
	}
	err = conn.QueryRow(ctx, `UPDATE contacts SET email = $1, name = $2 WHERE id = $3 RETURNING email, name, id`,
		c.Email, c.Name, c.Id).Scan(&c.Email, &c.Name, &c.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}
	return c, nil
}

func (r *ReportRepository) DeleteContact(ctx context.Context, contactId int64, userId int64) (bool, error) {
	conn := r.db.Db
	var exist bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`,
		contactId, userId).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exist {
		return false, ErrNotFound
	}
	_, err = conn.Exec(ctx, `DELETE FROM contacts WHERE id = $1`, contactId)
	if err != nil {
		return false, fmt.Errorf("failed to delete contact: %w", err)
	}
	return true, nil
}

func (r *ReportRepository) GetContactsByTemplate(ctx context.Context, templateId int64, userId int64) ([]*report.Contact, error) {
	conn := r.db.Db
	var contacts []*report.Contact
	rows, err := conn.Query(ctx, `SELECT c.id, c.email, c.name FROM contacts c 
    JOIN contacts_templates ct ON c.id = ct.contact_id 
    WHERE ct.template_id = $1 AND c.user_id = $2`, templateId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c report.Contact
		err = rows.Scan(&c.Id, &c.Email, &c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, &c)
	}
	return contacts, nil
}

func (r *ReportRepository) AddContactsToTemplate(ctx context.Context, templateId int64,
	contactIds []int64, userId int64) error {
	tx, err := r.db.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var templateExists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM templates WHERE id = $1 AND user_id = $2)`,
		templateId, userId).Scan(&templateExists)
	if err != nil {
		return fmt.Errorf("failed to verify template: %w", err)
	}
	if !templateExists {
		return ErrNotFound
	}

	var validContacts int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM contacts 
		WHERE id = ANY($1) AND user_id = $2`,
		contactIds, userId).Scan(&validContacts)
	if err != nil {
		return fmt.Errorf("failed to verify contacts: %w", err)
	}
	if validContacts != len(contactIds) {
		return ErrNotFound
	}
	for _, contactId := range contactIds {
		_, err := tx.Exec(ctx,
			`INSERT INTO contacts_templates (template_id, contact_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING`,
			templateId, contactId)
		if err != nil {
			return fmt.Errorf("failed to add contact to template: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *ReportRepository) GetContacts(ctx context.Context, id int64) ([]*report.Contact, error) {
	conn := r.db.Db
	var contacts []*report.Contact
	rows, err := conn.Query(ctx, `SELECT id, name, email FROM contacts WHERE user_id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var contact report.Contact
		if err := rows.Scan(&contact.Id, &contact.Name, &contact.Email); err != nil {
			return nil, fmt.Errorf("failed to scan contact row: %w", err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating contact rows: %w", err)
	}

	return contacts, nil
}
