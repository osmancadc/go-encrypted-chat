package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type CipherFactory interface {
	newCipher(key []byte) (cipher.Block, error)
	newGCM(block cipher.Block) (cipher.AEAD, error)
}

type Encryptor struct{}

func (d *Encryptor) newCipher(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

func (d *Encryptor) newGCM(block cipher.Block) (cipher.AEAD, error) {
	return cipher.NewGCM(block)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type AES struct {
	key []byte
}

func NewAES(key []byte) (*AES, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("the AES is invalid, it must have 16, 24 or 32 bytes")
	}
	return &AES{key: key}, nil
}

func (a *AES) GetKey() []byte {
	return a.key
}

func GenerateAES(size int, randReader Reader) (*AES, error) {
	key := make([]byte, size)

	if size != 16 && size != 24 && size != 32 {
		return nil, fmt.Errorf("error generating AES key, invalid size %d, must be 16, 28 or 32 bytes", size)
	}

	_, err := randReader.Read(key)
	if err != nil {
		return nil, err
	}

	return &AES{key: key}, nil
}

func generateNonce(randReader Reader) (nonce []byte, err error) {
	nonce = make([]byte, 12)

	if _, err := randReader.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func (a *AES) EncryptWithAESGCM(factory CipherFactory, randReader Reader, plaintext []byte) (ciphertext []byte, err error) {

	nonce, err := generateNonce(randReader)
	if err != nil {
		return
	}

	block, err := factory.newCipher(a.key)
	if err != nil {
		return
	}

	gcm, err := factory.newGCM(block)
	if err != nil {
		return
	}

	sealedMessage := gcm.Seal(nil, nonce, plaintext, nil)

	ciphertext = append(nonce, sealedMessage...)

	return
}

func (a *AES) DecryptWithAESGCM(factory CipherFactory, ciphertext []byte) (plaintext []byte, err error) {
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	block, err := factory.newCipher(a.key)
	if err != nil {
		return
	}

	gcm, err := factory.newGCM(block)
	if err != nil {
		return
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)

	return
}
