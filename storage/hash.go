package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hash string

func (h Hash) String() string {
	return string(h)
}

func createHash(password string) Hash {
	h := sha256.Sum256([]byte(password))
	return Hash(hex.EncodeToString(h[:]))
}
