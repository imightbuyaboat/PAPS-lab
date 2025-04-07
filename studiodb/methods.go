package studiodb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

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

func (db *DB) Insert(i Item) error {
	query := "INSERT INTO register (organization, city, phone) VALUES ($1, $2, $3)"

	_, err := db.Exec(query, i.Organization, i.City, i.Phone)
	return err
}

func (db *DB) SelectAll() ([]Item, error) {
	query := "SELECT * FROM register"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []Item{}
	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.Id, &i.Organization, &i.City, &i.Phone)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (db *DB) SelectAny(i Item) ([]Item, error) {
	query := "SELECT * FROM register where 1=1"
	var args []interface{}

	if i.Organization != "" {
		query += " and organization = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.Organization)
	}
	if i.City != "" {
		query += " and city = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.City)
	}
	if i.Phone != "" {
		query += " and phone = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.Phone)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []Item{}
	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.Id, &i.Organization, &i.City, &i.Phone)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (db *DB) Delete(id int) error {
	_, err := db.Exec("DELETE FROM register WHERE id = $1", id)
	return err
}

/*CREATE TABLE register (
    id SERIAL PRIMARY KEY,
    organization VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL
);*/
