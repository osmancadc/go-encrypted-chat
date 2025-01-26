package model

import (
	"crypto/rand"

	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
)

func EncryptMessage(plaintext []byte, aesInstance crypto.AES) (ciphertext []byte, err error) {

	cipherFactory := crypto.Encryptor{}
	ciphertext, err = aesInstance.EncryptWithAESGCM(&cipherFactory, rand.Reader, plaintext)

	return
}

func DecryptMessage(ciphertext []byte, aesInstance crypto.AES) (plaintext []byte, err error) {

	cipherFactory := crypto.Encryptor{}
	plaintext, err = aesInstance.DecryptWithAESGCM(&cipherFactory, ciphertext)

	return
}

func EncryptRSA(plaintext []byte, rsaInstance crypto.RSA) (ciphertext []byte, err error) {
	ciphertext, err = rsaInstance.EncryptMessage(plaintext)

	return
}

func DecryptRSA(ciphertext []byte, rsaInstance crypto.RSA) (plaintext []byte, err error) {
	plaintext, err = rsaInstance.DecryptMessage(ciphertext)

	return
}
