package postgres

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5435"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"reportDB"`
}

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DB struct {
	Db DBInterface
}

func New(config Config) (*DB, error) {
	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &DB{Db: db}, nil
}
