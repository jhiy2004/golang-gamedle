package game

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connectable interface {
	GetConnection() *websocket.Conn
	Connect(conn *websocket.Conn)
	Disconnect()
}

type PlayerState struct {
	Question  int
	Connected bool
}

type Player interface {
	Send(msg *Message) error
	Receive() (*Message, error)
	GetName() string
	GetState() *PlayerState

	IsReady() bool
	ToggleReady() bool

	IsRetry() bool
	ToggleRetry() bool
	Reset()
}

type WSPlayer struct {
	Conn   *websocket.Conn
	SendCh chan []byte
	Name   string
	Ready  bool
	Retry  bool
	Mu     *sync.Mutex
	PlayerState
}

type HostPlayer struct {
	Name    string
	Ready   bool
	Retry   bool
	Channel chan []byte
	PlayerState
}

func (p *WSPlayer) StartWriter() {
	go func() {
		for msg := range p.SendCh {
			err := p.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}()
}

func (p *WSPlayer) GetConnection() *websocket.Conn {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return p.Conn
}

func (p *WSPlayer) Connect(conn *websocket.Conn) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Conn = conn
	p.Connected = true
	p.SendCh = make(chan []byte, 16)

	p.StartWriter()
}

func (p *WSPlayer) Disconnect() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Conn = nil
	p.Connected = false
	close(p.SendCh)
}

func (p *WSPlayer) GetState() *PlayerState {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return &p.PlayerState
}

func (p *WSPlayer) Send(msg *Message) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	p.SendCh <- content
	return nil
}

func (p *WSPlayer) Receive() (*Message, error) {
	p.Mu.Lock()

	if p.Conn == nil {
		p.Mu.Unlock()
		return nil, errors.New("no connection")
	}

	conn := p.Conn
	p.Mu.Unlock()

	_, content, err := conn.ReadMessage()
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
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return p.Name
}

func (p *WSPlayer) ToggleReady() bool {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Ready = !p.Ready

	return p.Ready
}

func (p *WSPlayer) IsReady() bool {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return p.Ready
}

func (p *WSPlayer) ToggleRetry() bool {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Retry = !p.Retry

	return p.Retry
}

func (p *WSPlayer) IsRetry() bool {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return p.Retry
}

func (p *WSPlayer) Reset() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Ready = false
	p.Retry = false
	p.PlayerState.Question = 0
}

func (p *HostPlayer) GetState() *PlayerState {
	return &p.PlayerState
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
	p.PlayerState.Question = 0
}
