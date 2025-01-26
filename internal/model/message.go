package model

import (
	"encoding/json"
	"fmt"
)

type WebsocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (m *WebsocketMessage) Marshal() ([]byte, error) {
	messageBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return messageBytes, nil
}

func (m *WebsocketMessage) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &m)

	return err
}

type PublicKeyExchangePayload struct {
	PublicKey      []byte `json:"publicKey"`
	NeedsPublicKey bool   `json:"needPublicKey"`
	UserID         string `json:"userID"`
}

func (m *PublicKeyExchangePayload) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &m)

	return err
}

type TextMessagePayload struct {
	Content  string `json:"content"`
	SenderID string `json:"senderID"`
	GroupID  string `json:"groupID"`
}

func (m *TextMessagePayload) Marshal() ([]byte, error) {
	messageBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return messageBytes, nil
}

func (m *TextMessagePayload) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &m)

	return err
}

type UsernamePayload struct {
	Username string `json:"Username"`
}

func (m *UsernamePayload) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &m)

	return err
}

type IncomingMessage struct {
	Message TextMessagePayload
}

func (m IncomingMessage) String() string {
	return fmt.Sprintf("%s: %s", m.Message.SenderID, m.Message.Content)
}
