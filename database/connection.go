package database

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type ConnectionConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	SSLMode  string `json:"SSLMode" validate:"required"`
}

func (d *ConnectionConfig) AsDSN() string {
	if d.Port == 0 {
		d.Port = 5432
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/dynamic-forms?sslmode=%s", d.Username, d.Password, d.Host, d.Port, d.SSLMode)
}

func CreateDatabaseConnection(config ConnectionConfig) (*bun.DB, error) {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.AsDSN())))

	db := bun.NewDB(sqlDB, pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
