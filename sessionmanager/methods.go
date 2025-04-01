package sessionmanager

import (
	"math/rand"
	"sync"
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	id := SessionID{RandStringRunes(sessionKeyLen)}
	sm.mu.Lock()
	sm.sessions[id] = in
	sm.mu.Unlock()
	return &id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if session, ok := sm.sessions[*in]; ok {
		return session
	}
	return nil
}

func (sm *SessionManager) Delete(in *SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
}
