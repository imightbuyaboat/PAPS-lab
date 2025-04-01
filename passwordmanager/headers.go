package passwordmanager

import (
	"fmt"
	"sync"
)

type User struct {
	Login    string
	Password string
}

type Hash struct {
	hash []byte
}

type PasswordManager struct {
	mu    sync.RWMutex
	users map[string]*Hash
}

type MyError struct {
	errorId int
}

func (me *MyError) Error() string {
	return fmt.Sprint(me.errorId)
}
