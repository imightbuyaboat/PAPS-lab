package studiodb

import (
	"database/sql"
)

type Item struct {
	Id           int
	Organization string
	Phone        string
}

type DB struct {
	*sql.DB
}
