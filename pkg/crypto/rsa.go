package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type RSA struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func GenerateRSA(size int) (*RSA, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	return &RSA{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
	}, nil
}

func (r *RSA) GetPublicKeyValue() (publicKey []byte, err error) {
	publicKey, err = x509.MarshalPKIXPublicKey(r.publicKey)

	return
}

func (r *RSA) EncryptMessage(plaintext []byte) (ciphertext []byte, err error) {
	ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, plaintext)

	return
}

func (r *RSA) DecryptMessage(ciphertext []byte) (plaintext []byte, err error) {
	plaintext, err = rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, ciphertext)

	return
}
