package chat

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
)

type User struct {
	ID        string
	username  string
	publicKey []byte
}

func NewUser(username string) (*User, error) {

	uuid := uuid.NewString()

	rsaInstance, err := crypto.GenerateRSA(2048)
	if err != nil {
		return nil, err
	}

	publicKey, err := rsaInstance.GetPublicKeyValue()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid,
		username:  username,
		publicKey: publicKey,
	}, nil
}

func (u *User) GetUserPublicKey() []byte {
	return u.publicKey
}

func (u *User) Marshal() ([]byte, error) {
	userBytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return userBytes, nil
}

func (u *User) Unmarshal(data []byte) (user User, err error) {

	err = json.Unmarshal(data, &user)

	return
}
