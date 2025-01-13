package main

import (
	"fmt"

	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
)

func main() {
	fmt.Println("----------------------------------------------- SECURE CHAT INITIALIZED -----------------------------------------------")
	aesInstance, err := crypto.GenerateAES(32)
	if err != nil {
		fmt.Println("Error generating AES key")
		return
	}

	plainMessage := "hello world"
	encryptedMessage, err := aesInstance.EncryptWithAESGCM([]byte(plainMessage))
	if err != nil {
		fmt.Println("Error encrypting message")
		return
	}

	decryptedMessage, err := aesInstance.DecryptWithAESGCM(encryptedMessage)
	if err != nil {
		fmt.Println("Error decrypting message")
	}

	fmt.Println(string(decryptedMessage))
}
