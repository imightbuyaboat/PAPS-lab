package passwordmanager

import (
	"PAPS-LAB/studiodb"
)

type User struct {
	Login    string
	Password string
}

type Info struct {
	Hash       string
	Priveleged bool
}

type PasswordManager struct {
	*studiodb.DB
}

/*
CREATE TABLE users (
    login TEXT PRIMARY KEY,
    hash TEXT NOT NULL,
    priveleged BOOLEAN NOT NULL DEFAULT FALSE
);
*/
