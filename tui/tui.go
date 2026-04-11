package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/jhiy2004/golang-gamedle/server/game"
)

type StartMsg struct {
	Msg game.StartMsg
}

type LobbyMsg struct {
	Msg game.LobbyMsg
}

type NotifyMsg struct {
	Text string
}

type StateMsg struct {
	State game.StateMsg
}

type GameState struct {
	Question   string
	Player     string
	State      game.RoomState
	Winner     string
	MaxPlayers int
	MinPlayers int
}

type Model struct {
	Width                 int
	Height                int
	Notifications         []string
	State                 GameState
	Quitting              bool
	Viewport              viewport.Model
	SenderStyle           lipgloss.Style
	Textinput             textinput.Model
	PromptBoxStyle        lipgloss.Style
	NotificationsBoxStyle lipgloss.Style
	QuestionBoxStyle      lipgloss.Style

	SendMessageCallback func(*game.Message) error
	QuitCallback        func() error
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var invalidInput string

	switch msg := msg.(type) {
	case StartMsg:
		msgContent := msg.Msg

		m.State.Player = msgContent.PlayerName
		m.State.MinPlayers = msgContent.MinPlayers
		m.State.MaxPlayers = msgContent.MaxPlayers

	case LobbyMsg:
		msgContent := msg.Msg

		m.Notifications = append(
			m.Notifications,
			fmt.Sprintf("%d/%d Players in the room", msgContent.CurrPlayers, m.State.MaxPlayers),
		)

		m.Notifications = append(
			m.Notifications,
			fmt.Sprintf("%d/%d Players ready", msgContent.ReadyPlayers, msgContent.CurrPlayers),
		)

		m.Viewport.SetContent(
			lipgloss.NewStyle().
				Width(m.Viewport.Width()).
				Render(strings.Join(m.Notifications, "\n")),
		)
		m.Viewport.GotoBottom()

	case NotifyMsg:
		m.Notifications = append(m.Notifications, msg.Text)
		m.Viewport.SetContent(
			lipgloss.NewStyle().
				Width(m.Viewport.Width()).
				Render(strings.Join(m.Notifications, "\n")),
		)
		m.Viewport.GotoBottom()

	case StateMsg:
		m.State.Question = msg.State.Question
		m.State.Winner = msg.State.Winner
		m.State.State = game.StringToRoomState(msg.State.State)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		pbStyle := lipgloss.NewStyle().
			Width(m.Width).
			Height(3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Zest).
			Align(lipgloss.Left, lipgloss.Bottom)

		qbStyle := lipgloss.NewStyle().
			Width(m.Width).
			Height(3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Cherry).
			Align(lipgloss.Left, lipgloss.Top)

		nbStyle := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height-pbStyle.GetHeight()-qbStyle.GetHeight()).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Thunder).
			Align(lipgloss.Left, lipgloss.Top)

		m.PromptBoxStyle = pbStyle
		m.NotificationsBoxStyle = nbStyle
		m.QuestionBoxStyle = qbStyle

		m.Textinput.SetWidth(m.Width)

		m.Viewport.SetWidth(m.Width - 2)
		m.Viewport.SetHeight(m.Height - 2 - pbStyle.GetHeight() - qbStyle.GetHeight())
		m.Viewport.SetXOffset(1)
		m.Viewport.SetYOffset(1)

		if len(m.Notifications) > 0 {
			// Wrap content before setting it.
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width()).Render(strings.Join(m.Notifications, "\n")))
		}
		m.Viewport.GotoBottom()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.QuitCallback()
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			if m.State.State == game.Waiting && m.Textinput.Value() == "ready" {
				msg, err := game.NewReadyMsg()
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}

				err = m.SendMessageCallback(msg)
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}
			} else if m.State.State == game.Waiting && m.Textinput.Value() == "cancel" {
				msg, err := game.NewCancelMsg()
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}

				err = m.SendMessageCallback(msg)
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}
			} else if m.State.State == game.Playing {
				msg, err := game.NewGuessMsg(m.Textinput.Value())
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}

				err = m.SendMessageCallback(msg)
				if err != nil {
					m.Notifications = append(m.Notifications, m.SenderStyle.Render("Error: ")+err.Error())
				}
			} else if m.State.State == game.End {
				invalidInput = "You can't send more messages"
			} else {
				invalidInput = "Invalid input"
			}

			m.Notifications = append(m.Notifications, m.SenderStyle.Render("You: ")+m.Textinput.Value())

			if invalidInput != "" {
				errStyle := lipgloss.NewStyle().
					Foreground(charmtone.Cherry)
				m.Notifications = append(m.Notifications, m.SenderStyle.Render(errStyle.Render(invalidInput)))
			}

			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width()).Render(strings.Join(m.Notifications, "\n")))
			m.Textinput.Reset()
			m.Viewport.GotoBottom()
			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.Textinput, cmd = m.Textinput.Update(msg)
			return m, cmd
		}
	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.Textinput, cmd = m.Textinput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() tea.View {
	viewportView := m.Viewport.View()
	v := tea.NewView("")

	c := m.Textinput.Cursor()
	if c != nil {
		c.Y = m.Height - 2
		c.X += 1
	}
	v.Cursor = c
	v.AltScreen = true

	questionBox := lipgloss.NewLayer(
		m.QuestionBoxStyle.Render("Question: " + m.State.Question),
	).X(0).Y(0)

	notificationsBox := lipgloss.NewLayer(
		m.NotificationsBoxStyle.Render(viewportView),
	).X(0).Y(3)

	promptBox := lipgloss.NewLayer(
		m.PromptBoxStyle.Render(m.Textinput.View()),
	).X(0).Y(m.Height - 3)

	comp := lipgloss.NewCompositor(
		questionBox,
		notificationsBox,
		promptBox,
	)

	v.SetContent(comp.Render())
	return v
}

func InitModel(sendMessageCallback func(*game.Message) error, quitCallBack func() error) Model {
	ti := textinput.New()
	ti.Placeholder = "Send your guess..."
	ti.SetVirtualCursor(false)
	ti.Focus()

	ti.Prompt = "> "
	ti.CharLimit = 280

	ti.SetWidth(30)

	// Remove cursor line styling
	s := ti.Styles()
	ti.SetStyles(s)

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the game room!
Type a guess and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	return Model{
		Textinput:           ti,
		Notifications:       []string{},
		Viewport:            vp,
		SenderStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		State:               GameState{},
		SendMessageCallback: sendMessageCallback,
		QuitCallback:        quitCallBack,
	}
}
