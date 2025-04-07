package sessionmanager

import "sync"

type Session struct {
	Login      string
	Useragent  string
	Priveleged bool
}

type SessionID struct {
	ID string
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[SessionID]*Session
}

const sessionKeyLen = 10

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
