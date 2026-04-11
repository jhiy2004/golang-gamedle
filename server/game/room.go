package game

import (
	"encoding/json"
	"fmt"
	"github.com/jhiy2004/golang-gamedle/server/db"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type RoomState int

const (
	Waiting RoomState = iota
	Playing
	End
)

const (
	configFilename = "config.json"
)

type RoomConfig struct {
	MinPlayers     int           `json:"minPlayers"`
	MaxPlayers     int           `json:"maxPlayers"`
	QuestionsCount int           `json:"questionsCount"`
	TimeoutSec     time.Duration `json:"timeoutSec"`
}

type Room struct {
	MinPlayers, MaxPlayers, CurrPlayers int
	Players                             map[Player]bool
	Status                              RoomState
	Mu                                  *sync.Mutex
	Questions                           map[int]db.QuestionAnswersDTO
	QuestionsOrder                      []int
	ReadyPlayers                        int
	Winner                              Player

	MinReached *sync.Cond
	Ready      chan struct{}
}

func StringToRoomState(state string) RoomState {
	switch state {
	case "waiting":
		return Waiting
	case "playing":
		return Playing
	case "end":
		return End
	default:
		return Waiting
	}
}

func RoomStateToString(state RoomState) string {
	switch state {
	case Waiting:
		return "waiting"
	case Playing:
		return "playing"
	case End:
		return "end"
	default:
		return "waiting"
	}
}

func NewRoom(conf *RoomConfig) *Room {
	mutex := &sync.Mutex{}

	room := Room{
		MinPlayers:     conf.MinPlayers,
		MaxPlayers:     conf.MaxPlayers,
		CurrPlayers:    0,
		Status:         Waiting,
		Players:        make(map[Player]bool),
		Winner:         nil,
		Questions:      make(map[int]db.QuestionAnswersDTO, conf.QuestionsCount),
		QuestionsOrder: make([]int, 0, conf.QuestionsCount),
		ReadyPlayers:   0,
		MinReached:     sync.NewCond(mutex),
		Ready:          make(chan struct{}),
		Mu:             mutex,
	}

	return &room
}

func (r *Room) GetStatus() RoomState {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	return r.Status
}

func (r *Room) PlayerReady(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	_, ok := r.Players[player]
	if !ok {
		log.Fatal("[ERROR] there should be an entry for all players")
	}

	if !player.IsReady() {
		player.ToggleReady()
		r.ReadyPlayers++
		if r.ReadyPlayers == r.CurrPlayers {
			r.SignalReady()
		}
		return true
	}

	return false
}

func (r *Room) PlayerCancel(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.ReadyPlayers > 0 && player.IsReady() {
		r.ReadyPlayers--
		player.ToggleReady()
		return true
	}

	return false
}

func (r *Room) WaitReady() <-chan struct{} {
	return r.Ready
}

func (r *Room) SignalReady() {
	close(r.Ready)
}

func (r *Room) WaitMinReached() {
	r.MinReached.L.Lock()
	defer r.MinReached.L.Unlock()

	for r.CurrPlayers < r.MinPlayers {
		r.MinReached.Wait()
	}
}

func (r *Room) SignalMinReached() {
	r.MinReached.Broadcast()
}

func (r *Room) Add(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.CurrPlayers >= r.MaxPlayers {
		return false
	}

	r.Players[player] = true
	r.CurrPlayers++
	r.SignalMinReached()

	if r.CurrPlayers == r.MaxPlayers {
		fmt.Println("Why?", r.CurrPlayers, r.MaxPlayers)
		r.SignalReady()
	}

	return true
}

func (r *Room) Remove(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Players[player] = false
	r.CurrPlayers--
	r.SignalMinReached()

	return true
}

func (r *Room) Broadcast(player Player, message *Message) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for p, ok := range r.Players {
		if p == player {
			continue
		}

		if ok {
			err := p.Send(message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (r *Room) ValidateAnswer(questionId int, answer string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	question, ok := r.Questions[questionId]
	if !ok {
		log.Fatal("Not ok questionId")
	}

	for _, ans := range question.Answers {
		if strings.ToUpper(ans) == strings.ToUpper(answer) {
			return true
		}
	}

	return false
}

func ReadConfig() *RoomConfig {
	defaultConfig := RoomConfig{
		MinPlayers:     2,
		MaxPlayers:     3,
		QuestionsCount: 5,
		TimeoutSec:     5 * time.Second,
	}

	file, err := os.Open(configFilename)
	defer file.Close()

	if err != nil {
		content, err := json.Marshal(&defaultConfig)
		if err != nil {
			log.Fatal("Failed to create the config file")
		}

		// rw-r--r--
		err = os.WriteFile(configFilename, content, 0644)
		if err != nil {
			log.Fatal("Failed to write into config file")
		}

		return &defaultConfig
	}

	conf := RoomConfig{}
	dec := json.NewDecoder(file)
	err = dec.Decode(&conf)
	if err != nil {
		fmt.Println(err)
		return &defaultConfig
	}

	return &conf
}
