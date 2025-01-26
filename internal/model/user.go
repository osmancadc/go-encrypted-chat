package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID       string
	Username string
}

func NewUser(userID, username string) (*User, error) {

	if userID == "" {
		userID = uuid.NewString()
	}

	return &User{
		ID:       userID,
		Username: username,
	}, nil
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
