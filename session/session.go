package session

import (
	"github.com/google/uuid"
)

type Session struct {
	Login      string
	Useragent  string
	Priveleged bool
}

type SessionID uuid.UUID

func (sid SessionID) String() string {
	return uuid.UUID(sid).String()
}

func NewSessionID() (SessionID, error) {
	id, err := uuid.NewRandom()
	return SessionID(id), err
}

func ParseSessionID(s string) (SessionID, error) {
	id, err := uuid.Parse(s)
	return SessionID(id), err
}
