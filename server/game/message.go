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
	Player   string `json:"player"`
	Winner   string `json:"winner"`
	State    string `json:"state"`
}

type GuessMsg struct {
	Answer string `json:"answer"`
}

type ReadyMsg struct{}

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
		Player:   player,
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
