package chat

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          string
	SenderID    string
	RecipientID string
	RoomID      string
	Content     []byte
	Timestamp   time.Time
	MessageType MessageType
}

type MessageType int

const (
	DirectMessage MessageType = iota
	RoomMessage
)

func (m *Message) GetSenderID() string {
	return m.SenderID
}

func (m *Message) GetRecipientID() string {
	return m.RecipientID
}

func (m *Message) GetContent() []byte {
	return m.Content
}

func (m *Message) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *Message) GetType() MessageType {
	return m.MessageType
}

func (m *Message) NewDirectMessage(senderID, recipientID string, content []byte) *Message {

	return &Message{
		ID:          uuid.NewString(),
		SenderID:    senderID,
		RecipientID: recipientID,
		MessageType: DirectMessage,
		Content:     content,
		Timestamp:   time.Now(),
	}
}

func (m *Message) NewRoomMessage(senderID, roomID string, content []byte) *Message {
	return &Message{
		ID:          uuid.NewString(),
		SenderID:    senderID,
		RoomID:      roomID,
		MessageType: RoomMessage,
		Content:     content,
		Timestamp:   time.Now(),
	}
}

func (m *Message) Marshal() ([]byte, error) {
	messageBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return messageBytes, nil
}

func (m *Message) Unmarshal(data []byte) (message Message, err error) {
	err = json.Unmarshal(data, &message)

	return
}
