package config

import (
	"log"
	"sync"

	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
)

type Config struct {
	mu            sync.RWMutex
	rsaInstance   *crypto.RSA
	PublicKeys    map[string][]byte
	SymmetricKeys map[string][]byte
}

var (
	once     sync.Once
	instance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		rsaInstance, err := crypto.GenerateRSA(2048)
		if err != nil {
			log.Fatal("error creating the initial configuration")
		}
		instance = &Config{
			PublicKeys:    map[string][]byte{},
			SymmetricKeys: map[string][]byte{},
			rsaInstance:   rsaInstance,
		}
	})

	return instance
}

func (c *Config) GetRsaInstance() *crypto.RSA {
	c.mu.RLock()
	defer c.mu.RLock()
	return c.rsaInstance
}

func (c *Config) AddPublicKey(userID string, publicKey []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.PublicKeys[userID] = publicKey
}

func (c *Config) GetPublicKey(userID string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.PublicKeys[userID]
}

func (c *Config) RemovePublicKey(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.PublicKeys, userID)
}

func (c *Config) AddSymmetricKey(userID string, symmetricKey []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.SymmetricKeys[userID] = symmetricKey
}

func (c *Config) GetSymmetricKey(userID string) []byte {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SymmetricKeys[userID]
}

func (c *Config) RemoveSymmetricKey(userID string) {
	c.mu.Lock()
	defer c.mu.Lock()

	delete(c.SymmetricKeys, userID)
}
