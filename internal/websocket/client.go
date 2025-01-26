package websocket

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/osmancadc/go-encrypted-chat/internal/model"
	"github.com/osmancadc/go-encrypted-chat/internal/view"
	"github.com/osmancadc/go-encrypted-chat/pkg/logger"
)

var log = logger.NewLogger("INFO")

type ClientHandler struct {
	Conn            *Connection
	program         *tea.Program
	externalMsgChan chan model.IncomingMessage
}

func NewClientHandler(conn *Connection) *ClientHandler {
	return &ClientHandler{
		Conn:            conn,
		externalMsgChan: make(chan model.IncomingMessage),
	}
}

func (h *ClientHandler) Run() {
	u := "ws://localhost:8080/ws"

	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket: %v\n", err)
	}
	h.Conn.SetConn(conn)
	h.Conn.SetChat()

	go h.readPump()
	go h.writePump()

	log.Info("Client connected to server")

	h.sendMessage(model.WebsocketMessage{
		Type: "usernameMessage",
		Payload: model.UsernamePayload{
			Username: h.Conn.User.Username,
		},
	})

	chatModel := view.InitialModel(h.Conn.GetConn(), h.Conn.User.Username)
	h.program = tea.NewProgram(chatModel)

	go func() {
		for msg := range chatModel.Send {
			websocketMsg := model.WebsocketMessage{
				Type: "textMessage",
				Payload: model.TextMessagePayload{
					Content:  msg.Content,
					SenderID: msg.SenderID,
				},
			}
			h.sendMessage(websocketMsg)
		}
	}()

	go func() {
		for msg := range h.externalMsgChan {
			h.program.Send(msg)
		}
	}()

	_, err = h.program.Run()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}

func (h *ClientHandler) readPump() {
	defer h.Conn.Close()
	log.Debug("Entered to readPump")
	for {
		_, message, err := h.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("Read error: %v\n", err)
			}
			break
		}
		h.handleMessage(message)
	}
}

func (h *ClientHandler) writePump() {
	defer h.Conn.Close()
	log.Debug("Entered to write Pump")

	for message := range h.Conn.GetSendChan() {
		err := h.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Errorf("Write error: %v\n", err)
			break
		}
		log.Debug("Message sended successfully")
	}
}

func (h *ClientHandler) handleMessage(message []byte) (err error) {
	var chatMessage model.WebsocketMessage
	err = json.Unmarshal(message, &chatMessage)
	if err != nil {
		log.Errorf("Failed to unmarshal message: %s", err.Error())
		return
	}

	var textMsg model.TextMessagePayload

	byteMsg, err := json.Marshal(chatMessage.Payload)
	if err != nil {
		return
	}

	err = textMsg.Unmarshal(byteMsg)
	if err != nil {
		return
	}

	incomingMsg := model.IncomingMessage{Message: textMsg}
	h.externalMsgChan <- incomingMsg
	return nil
}

func (h *ClientHandler) sendMessage(msg model.WebsocketMessage) (err error) {
	log.Debug("Entered to send message")
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing: %v\n", err.Error())
		return
	}
	h.Conn.GetSendChan() <- msgBytes

	return
}
