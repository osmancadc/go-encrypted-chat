package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type AES struct {
	key []byte
}

func GenerateAES(size int) (*AES, error) {

	key := make([]byte, size)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return &AES{
		key: key,
	}, nil
}

func generateNonce() (nonce []byte, err error) {
	nonce = make([]byte, 12)

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func (a *AES) EncryptWithAESGCM(plaintext []byte) (ciphertext []byte, err error) {

	nonce, err := generateNonce()
	if err != nil {
		return
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	sealedMessage := gcm.Seal(nil, nonce, plaintext, nil)

	ciphertext = append(nonce, sealedMessage...)

	return
}

func (a *AES) DecryptWithAESGCM(ciphertext []byte) (plaintext []byte, err error) {
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)

	return
}
