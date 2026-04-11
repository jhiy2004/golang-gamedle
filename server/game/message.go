package game

import (
	"encoding/json"
)

type Message struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

type NotifyMsg struct {
	Text string `json:"text"`
}

type StateMsg struct {
	Question string `json:"question"`
	Winner   string `json:"winner"`
	State    string `json:"state"`
}

type GuessMsg struct {
	Answer string `json:"answer"`
}

type LobbyMsg struct {
	CurrPlayers  int `json:"currPlayers"`
	ReadyPlayers int `json:"readyPlayers"`
}

type StartMsg struct {
	MinPlayers int    `json:"minPlayers"`
	MaxPlayers int    `json:"maxPlayers"`
	PlayerName string `json:"playerName"`
}

type ReadyMsg struct{}

type CancelMsg struct{}

func NewLobbyMsg(currPlayers, readyPlayers int) (*Message, error) {
	msg := LobbyMsg{
		CurrPlayers:  currPlayers,
		ReadyPlayers: readyPlayers,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "lobby",
		Payload: content,
	}

	return &base, nil
}

func NewStartMsg(minPlayers, maxPlayers int, playerName string) (*Message, error) {
	msg := StartMsg{
		MinPlayers: minPlayers,
		MaxPlayers: maxPlayers,
		PlayerName: playerName,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "start",
		Payload: content,
	}

	return &base, nil
}

func NewGuessMsg(answer string) (*Message, error) {
	msg := GuessMsg{
		Answer: answer,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "guess",
		Payload: content,
	}

	return &base, nil
}

func NewReadyMsg() (*Message, error) {
	msg := ReadyMsg{}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "ready",
		Payload: content,
	}

	return &base, nil
}

func NewCancelMsg() (*Message, error) {
	msg := ReadyMsg{}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "cancel",
		Payload: content,
	}

	return &base, nil
}

func NewNotifyMsg(text string) (*Message, error) {
	msg := NotifyMsg{
		Text: text,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "notify",
		Payload: content,
	}

	return &base, nil
}

func NewStateMsg(question, player, winner, state string) (*Message, error) {
	msg := StateMsg{
		Question: question,
		Winner:   winner,
		State:    state,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "state",
		Payload: content,
	}

	return &base, nil
}
