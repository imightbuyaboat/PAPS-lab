package sessionmanager

import (
	"sync"

	"github.com/google/uuid"
)

func NewSessionID() (SessionID, error) {
	id, err := uuid.NewRandom()
	return SessionID(id), err
}

func ParseSessionID(s string) (SessionID, error) {
	id, err := uuid.Parse(s)
	return SessionID(id), err
}

func (sid SessionID) String() string {
	return uuid.UUID(sid).String()
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(s *Session) (*SessionID, error) {
	id, err := NewSessionID()
	if err != nil {
		return nil, err
	}

	sm.mu.Lock()
	sm.sessions[id] = s
	sm.mu.Unlock()
	return &id, nil
}

func (sm *SessionManager) Check(id SessionID) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[id]
}

func (sm *SessionManager) Delete(id SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, id)
}
