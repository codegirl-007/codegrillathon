package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"os"
)

type Database interface {
	Query(query string, args ...interface{}) (*sqlx.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}

type DBClient struct {
	db *sqlx.DB
}

func NewDbClient() (*DBClient, error) {
	dsn := os.Getenv("DATABASE_URL") + "?authToken=" + os.Getenv("DATABASE_TOKEN")
	db, err := sqlx.Open("libsql", dsn)
	if err != nil {
		fmt.Printf("error connecting to database: %v", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &DBClient{db: db}, nil
}

func (c *DBClient) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.db.Queryx(query, args...)
}

func (c *DBClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *DBClient) Close() error {
	return c.db.Close()
}
