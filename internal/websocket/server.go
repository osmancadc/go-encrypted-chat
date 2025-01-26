package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/osmancadc/go-encrypted-chat/internal/model"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[string]*ServerHandler)
	clientsMu sync.Mutex
)

type ServerHandler struct {
	Conn *Connection
}

func NewServerHandler(conn *Connection) *ServerHandler {
	return &ServerHandler{Conn: conn}
}

func ServeWs() {
	port := ":8080"
	http.HandleFunc("/ws", handleConnections)
	log.Infof("Server started on %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v\n", err)
		return
	}
	log.Debug("Connection upgraded successfully")

	user := setNewUser(conn)

	clientConnection := NewConnection(user)
	log.Debugf("New connection created with Username: %s\n", clientConnection.User.Username)

	log.Debugf("New connection created with ID: %s\n", clientConnection.User.Username)

	clientConnection.SetConn(conn)

	handler := NewServerHandler(clientConnection)
	log.Debug("New ServerHandler created")

	clientsMu.Lock()
	clients[clientConnection.ID] = handler
	clientsMu.Unlock()
	log.Debug("Client added to clients map")

	go handler.Run()
	log.Debug("handler.Run() called")

	log.Debug("handleConnections finished")
}

func setNewUser(conn *websocket.Conn) model.User {
	var msg model.WebsocketMessage
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var usernameMsg model.UsernamePayload

	usernameByte, err := json.Marshal(msg.Payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	usernameMsg.Unmarshal(usernameByte)

	log.Debug(usernameMsg.Username)
	return model.User{
		Username: usernameMsg.Username,
	}
}

func (h *ServerHandler) Run() {
	go h.readPump()
	go h.writePump()
}

func (h *ServerHandler) readPump() {
	defer func() {
		h.Conn.Close()
		clientsMu.Lock()
		delete(clients, h.Conn.ID)
		clientsMu.Unlock()
		log.Infof("Client %s disconnected\n", h.Conn.ID)
	}()

	for {
		_, message, err := h.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Errorf("Reading error: %v\n", err)
			}
			break
		}

		log.Debugf("Message received from client %s\n", h.Conn.User.Username)

		clientsMu.Lock()
		for _, client := range clients {
			if client.Conn.ID != h.Conn.ID {
				client.Conn.GetSendChan() <- message
			}
		}
		clientsMu.Unlock()
	}
}

func (h *ServerHandler) writePump() {
	defer h.Conn.Close()
	for message := range h.Conn.GetSendChan() {
		err := h.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Errorf("write error: %v\n", err)
			break
		}
		log.Debugf("Message sent to client %s\n", h.Conn.ID)
	}
}

func (h *ServerHandler) handleMessage(message []byte) error {
	log.Debugf("Server received message: %s\n", message)
	return nil
}
