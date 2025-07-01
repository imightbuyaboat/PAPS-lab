package handler

import (
	"papslab/session"
)

type SessionManager interface {
	Create(s *session.Session) (*session.SessionID, error)
	Check(id session.SessionID) (*session.Session, error)
	Delete(id session.SessionID) error
}
