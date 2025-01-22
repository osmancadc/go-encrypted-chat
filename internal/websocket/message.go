package websocket

import (
	"encoding/json"

	"github.com/osmancadc/go-encrypted-chat/internal/chat"
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

func (m *WebsocketMessage) Unmarshal(data []byte) (message WebsocketMessage, err error) {
	err = json.Unmarshal(data, &message)

	return
}

type PublicKeyExchangePayload struct {
	PublicKey      []byte `json:"publicKey"`
	NeedsPublicKey bool   `json:"needPublicKey"`
	UserID         string `json:"userID"`
}

type InviteToGroupPayload struct {
	GroupName    string    `json:"groupName"`
	InviteeID    string    `json:"inviteeID"`
	EncryptedKey []byte    `json:"encryptedKey"`
	GroupID      string    `json:"groupID"`
	InviterUser  chat.User `json:"inviterUser"`
}

type AcceptInvitePayload struct {
	GroupID      string    `json:"groupID"`
	GroupName    string    `json:"groupName"`
	EncryptedKey []byte    `json:"encryptedKey"`
	InviterUser  chat.User `json:"inviterUser"`
}

type TextMessagePayload struct {
	Content  string `json:"content"`
	SenderID string `json:"senderID"`
	GroupID  string `json:"groupID"`
}
