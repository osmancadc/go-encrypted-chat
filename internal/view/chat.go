package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/osmancadc/go-encrypted-chat/internal/model"
)

const gap = "\n\n"

type (
	errMsg error
)

type ChatModel struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	receiverStyle lipgloss.Style
	err           error
	conn          *websocket.Conn
	Username      string
	Send          chan model.TextMessagePayload
}

func InitialModel(conn *websocket.Conn, username string) ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to ShadowChat! Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#60d300")),
		receiverStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#22a5ff")),
		err:           nil,
		conn:          conn,
		Username:      username,
		Send:          make(chan model.TextMessagePayload),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()

			m.Send <- model.TextMessagePayload{Content: m.textarea.Value(), SenderID: m.Username}
			m.textarea.Reset()
		}
	case model.IncomingMessage:
		newModel := m
		sender := fmt.Sprintf("%s: ", msg.Message.SenderID)
		newModel.messages = append(newModel.messages, newModel.receiverStyle.Render(sender)+msg.Message.Content)
		newModel.viewport.SetContent(lipgloss.NewStyle().Width(newModel.viewport.Width).Render(strings.Join(newModel.messages, "\n")))
		newModel.viewport.GotoBottom()
		return newModel, nil
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
