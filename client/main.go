package main

import (
	"encoding/json"
	"examples/clidle/server/game"
	"examples/clidle/tui"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gorilla/websocket"
)

func quitClosure(conn *websocket.Conn) func() error {
	return func() error {
		err := conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Ok"),
			time.Now().Add(1*time.Second),
		)

		if err != nil {
			return err
		}

		return conn.Close()
	}
}

func sendMessageClosure(conn *websocket.Conn) func(*game.Message) error {
	return func(message *game.Message) error {
		bytes, err := json.Marshal(message)
		if err != nil {
			return err
		}
		return conn.WriteMessage(websocket.TextMessage, bytes)
	}
}

func main() {
	dialer := &websocket.Dialer{
		NetDial:           nil,
		NetDialContext:    nil,
		NetDialTLSContext: nil,
		Proxy:             nil,
		TLSClientConfig:   nil,
		HandshakeTimeout:  (1600 * time.Millisecond),
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		Subprotocols:      nil,
		EnableCompression: false,
		Jar:               nil,
	}

	header := http.Header{}
	conn, _, err := dialer.Dial("ws://localhost:8080/ws", header)
	if err != nil {
		fmt.Println(err)
		return
	}

	myModel := tui.InitModel(sendMessageClosure(conn), quitClosure(conn))
	p := tea.NewProgram(myModel)

	// Notification logic
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			response := game.Message{}
			err = json.Unmarshal(message, &response)
			if err != nil {
				p.Send(tui.NotifyMsg{Text: err.Error()})
				continue
			}

			switch response.Cmd {
			case "notify":
				msg := game.NotifyMsg{}
				err = json.Unmarshal(response.Payload, &msg)
				if err != nil {
					p.Send(tui.NotifyMsg{Text: err.Error()})
					continue
				}

				p.Send(tui.NotifyMsg{Text: msg.Text})
			case "state":
				msg := game.StateMsg{}
				err = json.Unmarshal(response.Payload, &msg)
				if err != nil {
					p.Send(tui.NotifyMsg{Text: err.Error()})
					continue
				}

				p.Send(tui.StateMsg{State: msg})
			default:
				myModel.Notifications = append(myModel.Notifications, string(message))
				p.Send(tui.NotifyMsg{Text: string(message)})
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
