package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	config "github.com/osmancadc/go-encrypted-chat/configs"
	"github.com/osmancadc/go-encrypted-chat/internal/chat"
	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
	"github.com/osmancadc/go-encrypted-chat/pkg/logger"
)

type Client struct {
	conn *websocket.Conn
	user *chat.User
}

var (
	encryptor = crypto.Encryptor{}
	log       = logger.NewLogger("CHAT")
)

func (c *Client) HandleConnection() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, err := c.readMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err.Error())
			}
			break
		}

		err = c.handleMessage(message)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
}

func (c *Client) readMessage() (message WebsocketMessage, err error) {
	var msg WebsocketMessage

	err = c.conn.ReadJSON(&msg)

	return
}

func (c *Client) writeMessage(_ WebsocketMessage) (err error) {
	return
}

func (c *Client) handleMessage(msg WebsocketMessage) (err error) {
	switch msg.Type {
	case "publicKeyExchange":
		var payload PublicKeyExchangePayload
		err := json.Unmarshal([]byte(msg.Payload.(string)), &payload)
		if err != nil {
			return err
		}
		return c.handlePublicKeyExchange(payload)
	case "textMessage":
		var payload TextMessagePayload
		err := json.Unmarshal([]byte(msg.Payload.(string)), &payload)
		if err != nil {
			return err
		}
		return c.handleTextMessage(payload)

	case "inviteToGroup":
		var payload InviteToGroupPayload
		err := json.Unmarshal([]byte(msg.Payload.(string)), &payload)
		if err != nil {
			return err
		}
		return c.handleInviteToGroup(payload)
	case "acceptInvite":
		var payload AcceptInvitePayload
		err := json.Unmarshal([]byte(msg.Payload.(string)), &payload)
		if err != nil {
			return err
		}
		return c.handleAcceptInvite(payload)
	default:
		log.Info("Tipo de mensaje desconocido")
		return nil
	}
}

func (c *Client) handlePublicKeyExchange(payload PublicKeyExchangePayload) (err error) {

	if payload.NeedsPublicKey {
		msg := WebsocketMessage{
			Type: "publicKeyExchange",
			Payload: PublicKeyExchangePayload{
				PublicKey:      c.user.GetUserPublicKey(),
				NeedsPublicKey: false,
				UserID:         c.user.ID,
			},
		}

		c.writeMessage(msg)
	}

	return
}

func (c *Client) handleTextMessage(payload TextMessagePayload) (err error) {
	key := config.GetSymmetricKey(payload.SenderID)

	aesInstance, err := crypto.NewAES(key)
	if err != nil {
		return
	}

	plaintext, err := aesInstance.DecryptWithAESGCM(&encryptor, []byte(payload.Content))
	if err != nil {
		return
	}

	log.Log(string(plaintext))
	return
}

func (c *Client) handleInviteToGroup(_ InviteToGroupPayload) (err error) {
	// TODO: Implement invite to groups in the chat
	log.Info("Group feature it's still not available, we're sorry.")
	return
}

func (c *Client) handleAcceptInvite(_ AcceptInvitePayload) (err error) {
	log.Info("Group feature it's still not available, we're sorry.")
	// TODO: Implement accept groups in the chat
	return
}
