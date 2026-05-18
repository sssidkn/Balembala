package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5435"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"reportDB"`

	MaxConns int32 `yaml:"POSTGRES_MAX_CONNS" env:"POSTGRES_MAX_CONNS" env-default:"10"`
	MinConns int32 `yaml:"POSTGRES_MIN_CONNS" env:"POSTGRES_MIN_CONNS" env-default:"2"`
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
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.Database,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	if config.MaxConns > 0 {
		poolCfg.MaxConns = config.MaxConns
	}
	if config.MinConns > 0 {
		poolCfg.MinConns = config.MinConns
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &DB{Db: pool}, nil
}

func (d *DB) Close() {
	if pool, ok := d.Db.(*pgxpool.Pool); ok {
		pool.Close()
	}
}
