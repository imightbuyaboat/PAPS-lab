package basic_types

import (
	"github.com/google/uuid"
)

type User struct {
	Login    string
	Password string
}

type UserInfo struct {
	Hash       string
	Priveleged bool
}

type Item struct {
	Id           int
	Organization string
	City         string
	Phone        string
}

type Session struct {
	Login      string
	Useragent  string
	Priveleged bool
}

type SessionID uuid.UUID

func (sid SessionID) String() string {
	return uuid.UUID(sid).String()
}
