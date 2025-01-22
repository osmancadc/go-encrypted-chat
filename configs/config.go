package config

import (
	"sync"
)

var (
	symmetricKeysMutex sync.RWMutex
	SymmetricKeys      = make(map[string][]byte)
)

func AddSymmetricKey(userID string, symmetricKey []byte) {
	symmetricKeysMutex.Lock()
	defer symmetricKeysMutex.Unlock()
	SymmetricKeys[userID] = symmetricKey
}

func GetSymmetricKey(userID string) []byte {
	symmetricKeysMutex.RLock()
	defer symmetricKeysMutex.RUnlock()
	return SymmetricKeys[userID]
}

func DeleteSymmetricKey(userID string) {
	symmetricKeysMutex.Lock()
	defer symmetricKeysMutex.Unlock()
	delete(SymmetricKeys, userID)
}
