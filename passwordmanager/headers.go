package passwordmanager

import (
	"sync"
)

type User struct {
	Login    string
	Password string
}

type Info struct {
	Hash       []byte
	Priveleged bool
}

type PasswordManager struct {
	mu    sync.RWMutex
	users map[string]*Info
}
