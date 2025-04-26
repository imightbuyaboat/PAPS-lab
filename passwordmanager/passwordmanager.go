package passwordmanager

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	bt "papslab/basic_types"
	"papslab/studiodb"
)

type PasswordManager struct {
	*studiodb.DB
}

func CreateHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func NewPasswordManager(db *studiodb.DB) *PasswordManager {
	return &PasswordManager{db}
}

func (pm *PasswordManager) Insert(in *bt.User) error {
	hash := CreateHash(in.Password)

	query := "INSERT INTO users (login, hash) VALUES ($1, $2)"

	_, err := pm.Exec(query, in.Login, hash)
	return err
}

func (pm *PasswordManager) Check(in *bt.User) (exists bool, isPriv bool, err error) {
	query := "SELECT hash, priveleged FROM users where login = $1"
	row := pm.QueryRow(query, in.Login)

	var hashFromDB string
	var isPrivileged bool

	err = row.Scan(&hashFromDB, &isPrivileged)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}

	hash := CreateHash(in.Password)
	if hashFromDB == hash {
		return true, isPrivileged, nil
	}
	return false, false, nil
}

func (pm *PasswordManager) IsLoginAvailable(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE login=$1)`
	err := pm.QueryRow(query, login).Scan(&exists)
	return exists, err
}
