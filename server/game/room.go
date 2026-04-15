package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jhiy2004/golang-gamedle/server/db"
)

type RoomState int

const (
	Waiting RoomState = iota
	Playing
	End
	Restarting
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
	Players                             map[string]Player
	Status                              RoomState
	Mu                                  *sync.Mutex
	Questions                           map[int]db.QuestionAnswersDTO
	QuestionsOrder                      []int
	ReadyPlayers                        int
	RetryPlayers                        int
	Winner                              Player

	MinReached *sync.Cond
	Ready      chan struct{}
	IsEnded    chan struct{}
	Retry      chan struct{}
}

func GenerateRoomUUID() string {
	return uuid.NewString()
}

func GeneratePlayerUUID() string {
	return uuid.NewString()
}

func (r *Room) TryRestart() bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.Status != End {
		return false
	}

	r.Status = Restarting
	return true
}

func (r *Room) EndGame(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.Status == End {
		return false
	}

	r.Winner = player
	r.Status = End
	r.SignalIsEnded()

	return true
}

func (r *Room) Reset() {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Questions = make(map[int]db.QuestionAnswersDTO)
	r.QuestionsOrder = make([]int, 0)

	r.Status = Waiting
	r.ReadyPlayers = 0
	r.Ready = make(chan struct{})
	r.RetryPlayers = 0
	r.Retry = make(chan struct{})
	r.IsEnded = make(chan struct{})
	r.Winner = &WSPlayer{}

	for _, player := range r.Players {
		player.Reset()
	}
}

func (r *Room) Start(mydb *sql.DB, rng *rand.Rand, qtd int) error {
	ids, err := db.GetQuestionsIds(mydb)

	if qtd > len(ids) {
		return errors.New("qtd is greater than quantity of avaiable questions")
	}

	if err != nil {
		return err
	}

	perm := rng.Perm(len(ids))
	choices := perm[:qtd]

	for i, c := range choices {
		choices[i] = ids[c]
	}

	res, err := db.GetQuestions(mydb, choices)
	if err != nil {
		return err
	}

	for _, q := range res {
		r.Questions[q.Id] = q
	}
	r.QuestionsOrder = append(r.QuestionsOrder, choices...)

	return nil
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
		Players:        make(map[string]Player),
		Winner:         nil,
		Questions:      make(map[int]db.QuestionAnswersDTO, conf.QuestionsCount),
		QuestionsOrder: make([]int, 0, conf.QuestionsCount),
		ReadyPlayers:   0,
		RetryPlayers:   0,
		MinReached:     sync.NewCond(mutex),
		Ready:          make(chan struct{}),
		Retry:          make(chan struct{}),
		IsEnded:        make(chan struct{}),
		Mu:             mutex,
	}

	return &room
}

func (r *Room) GetStatus() RoomState {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	return r.Status
}

func (r *Room) PlayerRetry(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if !player.IsRetry() {
		player.ToggleRetry()
		r.RetryPlayers++

		if r.RetryPlayers == r.CurrPlayers {
			r.SignalRetry()
		}

		return true
	}

	return false
}

func (r *Room) PlayerCancelRetry(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.RetryPlayers > 0 && player.IsRetry() {
		player.ToggleRetry()
		r.RetryPlayers--

		return true
	}

	return false
}

func (r *Room) PlayerReady(player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

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

func (r *Room) WaitIsEnded() <-chan struct{} {
	return r.IsEnded
}

func (r *Room) SignalIsEnded() {
	close(r.IsEnded)
}

func (r *Room) WaitReady() <-chan struct{} {
	return r.Ready
}

func (r *Room) SignalReady() {
	close(r.Ready)
}

func (r *Room) WaitRetry() <-chan struct{} {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	return r.Retry
}

func (r *Room) SignalRetry() {
	close(r.Retry)
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

func (r *Room) PlayerExists(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	_, ok := r.Players[id]
	return ok
}

func (r *Room) GetPlayer(id string) Player {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	return r.Players[id]
}

func (r *Room) Add(id string, player Player) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.CurrPlayers >= r.MaxPlayers {
		return false
	}

	log.Printf("Added player with id: %s\n", id)
	r.Players[id] = player
	r.CurrPlayers++
	r.SignalMinReached()

	if r.CurrPlayers == r.MaxPlayers {
		log.Println("Curr players equals Max players")
		r.SignalReady()
	}

	return true
}

func (r *Room) Remove(id string) bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	delete(r.Players, id)

	r.CurrPlayers--
	r.SignalMinReached()

	return true
}

func (r *Room) Broadcast(player Player, message *Message) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for _, p := range r.Players {
		if p == player {
			continue
		}

		err := p.Send(message)
		if err != nil {
			fmt.Println(err)
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
