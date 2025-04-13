package sessionmanager

import (
	"sync"

	"github.com/google/uuid"
)

type Session struct {
	Login      string
	Useragent  string
	Priveleged bool
}

type SessionID uuid.UUID

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[SessionID]*Session
}
