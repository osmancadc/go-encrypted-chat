package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/osmancadc/go-encrypted-chat/config"
	"github.com/osmancadc/go-encrypted-chat/internal/chat"
	"github.com/osmancadc/go-encrypted-chat/pkg/crypto"
	"github.com/osmancadc/go-encrypted-chat/pkg/logger"
)

type Client struct {
	conn *websocket.Conn
	user *chat.User
}

var (
	encryptor  = crypto.Encryptor{}
	log        = logger.NewLogger("CHAT")
	keysConfig = config.GetConfig()
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

func (c *Client) writeMessage(msg WebsocketMessage) (err error) {
	writter, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	encodedMessage, err := msg.Marshal()
	if err != nil {
		return
	}

	writter.Write(encodedMessage)

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

func (c *Client) handlePublicKeyExchange(payload PublicKeyExchangePayload) error {

	if payload.NeedsPublicKey {
		rsaInstance := keysConfig.GetRsaInstance()
		publicKey, err := rsaInstance.GetPublicKeyValue()
		if err != nil {
			return err
		}

		msg := WebsocketMessage{
			Type: "publicKeyExchange",
			Payload: PublicKeyExchangePayload{
				PublicKey:      publicKey,
				NeedsPublicKey: false,
				UserID:         c.user.ID,
			},
		}

		c.writeMessage(msg)
	}

	return nil
}

func (c *Client) handleTextMessage(payload TextMessagePayload) (err error) {
	key := keysConfig.GetSymmetricKey(payload.SenderID)

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
