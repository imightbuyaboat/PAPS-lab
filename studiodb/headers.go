package studiodb

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}
