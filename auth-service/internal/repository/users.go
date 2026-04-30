package repository

import (
	"auth/pkg/db/postgres"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

import (
	"auth/internal/models"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("repository.Register: user already exists")
	ErrInvalidPassword   = errors.New("repository.Login: invalid password")
	ErrUserNotFound      = errors.New("repository.Login: user not found")
)

type Repository interface {
	Register(ctx context.Context, r *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, request *models.LoginRequest) (*models.User, error)
}

type UserRepository struct {
	db *postgres.DB
}

func NewUsersRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Register(ctx context.Context, r *models.RegisterRequest) (*models.User, error) {
	conn := repo.db.Db
	var userID int
	err := conn.QueryRow(ctx, `
    INSERT INTO users (email, username, password_hash) 
    VALUES ($1, $2, $3)
    RETURNING id`,
		r.Email, r.Username, r.Password,
	).Scan(&userID)

	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("repository.Register: %w", err)
	}

	user := &models.User{
		UserID:   int64(userID),
		Email:    r.Email,
		Username: r.Username,
	}
	return user, nil
}

func (repo *UserRepository) Login(ctx context.Context, request *models.LoginRequest) (*models.User, error) {
	conn := repo.db.Db
	email := request.Email
	password := request.Password
	var passwordHash string
	var id int64
	var username string
	err := conn.QueryRow(ctx, "SELECT username, password_hash, id FROM users WHERE email = $1",
		email).Scan(&username, &passwordHash, &id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}
	user := &models.User{
		UserID:   id,
		Email:    email,
		Username: username,
	}
	return user, nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
