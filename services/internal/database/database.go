package database

import (
	"context"
	"database/sql"

	"github.com/badrchoubai/services/internal/config"
)

var _ IDatabase = (*Database)(nil)

type Database struct {
	db *sql.DB
}

type IDatabase interface {
	Close() error
	DB() *sql.DB
	Ping(ctx context.Context) error
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func NewDatabase(cfg *config.AppConfig) (*Database, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func connect(cfg *config.AppConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DbConnectionString())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns())
	db.SetMaxIdleConns(cfg.MaxIdleConns())
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime())

	return db, nil
}
