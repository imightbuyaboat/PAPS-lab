package passwordmanager

import (
	"PAPS-LAB/studiodb"
	"crypto/sha256"
	"encoding/hex"
)

func CreateHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func NewPasswordManager(db *studiodb.DB) *PasswordManager {
	return &PasswordManager{db}
}

func (pm *PasswordManager) Insert(in *User) error {
	hash := CreateHash(in.Password)

	query := "INSERT INTO users (login, hash) VALUES ($1, $2)"

	_, err := pm.Exec(query, in.Login, hash)
	return err
}

func (pm *PasswordManager) Check(in *User) (int, bool, error) {
	query := "SELECT hash, priveleged FROM users where login = $1"
	rows, err := pm.Query(query, in.Login)
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()

	users := []Info{}
	for rows.Next() {
		i := Info{}
		err := rows.Scan(&i.Hash, &i.Priveleged)
		if err != nil {
			return 0, false, err
		}
		users = append(users, i)
	}

	if len(users) == 0 {
		return 2, false, err
	}

	hash := CreateHash(in.Password)
	if users[0].Hash == hash {
		return 0, users[0].Priveleged, err
	}
	return 1, false, err
}

func (pm *PasswordManager) CheckAvailableLogin(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE login=$1)`
	err := pm.QueryRow(query, login).Scan(&exists)
	return exists, err
}
