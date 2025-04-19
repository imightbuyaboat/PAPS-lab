package studiodb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	err := godotenv.Load() // Загружаем .env
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

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}
