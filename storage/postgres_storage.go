package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	*sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	host := os.Getenv("SQL_HOST")
	port := os.Getenv("SQL_PORT")
	dbname := os.Getenv("SQL_DB")
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresStorage{db}, nil
}
