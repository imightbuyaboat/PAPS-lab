package passwordmanager

import (
	"bytes"
	"crypto/sha256"
	"sync"
)

func CreateHash(password string) Hash {
	h := sha256.New()
	h.Write([]byte(password))
	return Hash{hash: (h.Sum(nil))}
}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		mu:    sync.RWMutex{},
		users: map[string]*Hash{},
	}
}

func (pm *PasswordManager) Create(in *User) error {
	hash := CreateHash(in.Password)
	pm.mu.Lock()
	pm.users[in.Login] = &hash
	pm.mu.Unlock()
	return nil
}

func (pm *PasswordManager) Check(in *User) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if userHash, ok := pm.users[in.Login]; ok {
		hash := CreateHash(in.Password)
		if bytes.Equal(userHash.hash, hash.hash) {
			return 0
		}
		return 1
	}
	return 2
}

func (pm *PasswordManager) CheckAvailableLogin(login string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if _, ok := pm.users[login]; ok {
		return 1
	}
	return 0
}
