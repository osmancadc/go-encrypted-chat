package main

import (
	"fmt"

	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
)

func main() {
	fmt.Println("----------------------------------------------- SECURE CHAT INITIALIZED -----------------------------------------------")

	rsaInstance, _ := crypto.GenerateRSA(2048)
	aesInstance, _ := crypto.GenerateAES(32)

	encrypted, _ := rsaInstance.EncryptMessage(aesInstance.GetKey())
	decrypted, _ := rsaInstance.DecryptMessage(encrypted)
	fmt.Println("===========================================DECRYPTED===========================================")
	fmt.Println(decrypted)

}
