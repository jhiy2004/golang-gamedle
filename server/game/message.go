package game

import (
	"encoding/json"
)

type Message struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

type GuessResponseMsg struct {
	Correct bool   `json:"correct"`
	Text    string `json:"text"`
}

type PlayerStatusMsg struct {
	Player   string `json:"player"`
	Progress int    `json:"progress"`
}

type NotifyMsg struct {
	Text string `json:"text"`
}

type StateMsg struct {
	Players  []string `json:"players"`
	Question string   `json:"question"`
	Winner   string   `json:"winner"`
	State    string   `json:"state"`
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
	PlayerId   string `json:"playerId"`
}

type PostGameLobbyMsg struct {
	CurrPlayers  int `json:"currPlayers"`
	RetryPlayers int `json:"retryPlayers"`
}

type ReadyMsg struct{}

type CancelMsg struct{}

type RetryMsg struct{}

type CancelRetryMsg struct{}

type RestartMsg struct{}

func NewRestartMsg() (*Message, error) {
	msg := RestartMsg{}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "restart",
		Payload: content,
	}

	return &base, nil
}

func NewRetryMsg() (*Message, error) {
	msg := RetryMsg{}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "retry",
		Payload: content,
	}

	return &base, nil
}

func NewCancelRetryMsg() (*Message, error) {
	msg := CancelRetryMsg{}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "cancelRetry",
		Payload: content,
	}

	return &base, nil
}

func NewPostGameLobbyMsg(currPlayers, retryPlayers int) (*Message, error) {
	msg := PostGameLobbyMsg{
		CurrPlayers:  currPlayers,
		RetryPlayers: retryPlayers,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "postGameLobby",
		Payload: content,
	}

	return &base, nil
}

func NewGuessResponseMsg(correct bool, text string) (*Message, error) {
	msg := GuessResponseMsg{
		Correct: correct,
		Text:    text,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "guessResponse",
		Payload: content,
	}

	return &base, nil
}

func NewPlayerStatusMsg(player string, progress int) (*Message, error) {
	msg := PlayerStatusMsg{
		Player:   player,
		Progress: progress,
	}

	content, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	base := Message{
		Cmd:     "playerStatus",
		Payload: content,
	}

	return &base, nil
}

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

func NewStartMsg(minPlayers, maxPlayers int, playerName, playerId string) (*Message, error) {
	msg := StartMsg{
		MinPlayers: minPlayers,
		MaxPlayers: maxPlayers,
		PlayerName: playerName,
		PlayerId:   playerId,
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

func NewStateMsg(question, winner, state string, players []string) (*Message, error) {
	msg := StateMsg{
		Players:  players,
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
