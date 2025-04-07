package passwordmanager

import (
	"bytes"
	"crypto/sha256"
	"sync"
)

func CreateHash(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	return h.Sum(nil)
}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		mu:    sync.RWMutex{},
		users: map[string]*Info{"admin": {CreateHash("admin"), true}},
	}
}

func (pm *PasswordManager) Create(in *User) error {
	hash := CreateHash(in.Password)
	pm.mu.Lock()
	pm.users[in.Login] = &Info{Hash: hash, Priveleged: false}
	pm.mu.Unlock()
	return nil
}

func (pm *PasswordManager) Check(in *User) (int, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if userInfo, ok := pm.users[in.Login]; ok {
		hash := CreateHash(in.Password)
		if bytes.Equal(userInfo.Hash, hash) {
			return 0, userInfo.Priveleged
		}
		return 1, false
	}
	return 2, false
}

func (pm *PasswordManager) CheckAvailableLogin(login string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if _, ok := pm.users[login]; ok {
		return 1
	}
	return 0
}
