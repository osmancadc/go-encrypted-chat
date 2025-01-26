package websocket

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/osmancadc/go-encrypted-chat/internal/model"
	"github.com/osmancadc/go-encrypted-chat/internal/view"
)

type Connection struct {
	ID   string
	conn *websocket.Conn
	send chan []byte
	User model.User
	chat *tea.Program
}

func NewConnection(user model.User) *Connection {
	return &Connection{
		ID:   uuid.New().String(),
		send: make(chan []byte, 256),
		User: user,
	}
}

func (c *Connection) SetChat() {
	c.chat = tea.NewProgram(view.InitialModel(c.conn, c.User.Username))
}

func (c *Connection) GetChat() *tea.Program {
	return c.chat
}

func (c *Connection) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *Connection) GetConn() *websocket.Conn {
	return c.conn
}

func (c *Connection) GetSendChan() chan []byte {
	return c.send
}

func (c *Connection) WriteMessage(messageType int, data []byte) error {
	if c.conn != nil {
		return c.conn.WriteMessage(messageType, data)
	}
	return nil
}

func (c *Connection) Close() error {
	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return err
		}
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) ReadMessage() (messageType int, p []byte, err error) {
	if c.conn != nil {
		return c.conn.ReadMessage()
	}
	return 0, nil, nil
}

func (c *Connection) GetUser() model.User {
	return c.User
}
