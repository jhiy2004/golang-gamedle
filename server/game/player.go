package game

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Player interface {
	Send(msg *Message) error
	Receive() (*Message, error)
	GetName() string

	IsReady() bool
	ToggleReady() bool

	IsRetry() bool
	ToggleRetry() bool
	Reset()
}

type WSPlayer struct {
	Conn  *websocket.Conn
	Name  string
	Ready bool
	Retry bool
	Mu    *sync.Mutex
}

type HostPlayer struct {
	Name    string
	Ready   bool
	Retry   bool
	Channel chan []byte
}

func (p *WSPlayer) Send(msg *Message) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Conn.WriteMessage(websocket.TextMessage, content)
}

func (p *WSPlayer) Receive() (*Message, error) {
	_, content, err := p.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	msg := Message{}
	err = json.Unmarshal(content, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (p *WSPlayer) GetName() string {
	return p.Name
}

func (p *WSPlayer) ToggleReady() bool {
	p.Ready = !p.Ready

	return p.Ready
}

func (p *WSPlayer) IsReady() bool {
	return p.Ready
}

func (p *WSPlayer) ToggleRetry() bool {
	p.Retry = !p.Retry

	return p.Retry
}

func (p *WSPlayer) IsRetry() bool {
	return p.Retry
}

func (p *WSPlayer) Reset() {
	p.Ready = false
	p.Retry = false
}

func (p *HostPlayer) Send(msg *Message) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	p.Channel <- content

	return nil
}

func (p *HostPlayer) Receive() (*Message, error) {
	msg := <-p.Channel

	cmdMsg := Message{}
	err := json.Unmarshal(msg, &cmdMsg)
	if err != nil {
		return nil, err
	}

	return &cmdMsg, nil
}

func (p *HostPlayer) GetName() string {
	return p.Name
}

func (p *HostPlayer) ToggleReady() bool {
	p.Ready = !p.Ready

	return p.Ready
}

func (p *HostPlayer) IsReady() bool {
	return p.Ready
}

func (p *HostPlayer) ToggleRetry() bool {
	p.Retry = !p.Retry

	return p.Retry
}

func (p *HostPlayer) IsRetry() bool {
	return p.Retry
}

func (p *HostPlayer) Reset() {
	p.Ready = false
	p.Retry = false
}
