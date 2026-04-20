package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"sync"

	tea "charm.land/bubbletea/v2"

	"github.com/gorilla/websocket"

	"github.com/jhiy2004/golang-gamedle/server/db"
	"github.com/jhiy2004/golang-gamedle/server/game"
	"github.com/jhiy2004/golang-gamedle/tui"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RoomResponseDTO struct {
	Id string `json:"id"`
}

type ListRoomsResponseDTO struct {
	Rooms []RoomResponseDTO `json:"rooms"`
}

type CreateRoomResponseDTO struct {
	Id string `json:"id"`
}

func randomColor(used []string) string {
	colors := []string{
		"red",
		"blue",
		"green",
		"yellow",
		"orange",
		"purple",
		"pink",
		"brown",
		"black",
		"white",
		"gray",
		"cyan",
		"magenta",
		"lime",
		"maroon",
		"navy",
		"olive",
		"teal",
		"silver",
		"gold",
		"coral",
		"indigo",
		"violet",
		"turquoise",
		"beige",
		"khaki",
		"lavender",
		"salmon",
		"plum",
		"tan",
	}

	avaiable := slices.DeleteFunc(colors, func(c string) bool {
		return slices.Contains(used, c)
	})

	if len(avaiable) == 0 {
		return ""
	}

	i := rand.Intn(len(avaiable))

	return colors[i]
}

func handleConnection(room *game.Room, playerId string, mydb *sql.DB, rng *rand.Rand, qtd int, killRoom func()) {
	log.Println("Handle new connection")
	defer killRoom()

	player := room.GetPlayer(playerId)

	wsPlayer, ok := player.(*game.WSPlayer)
	if ok {
		conn := wsPlayer.GetConnection()
		defer conn.Close()

		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("Oh shit: %d %s\n", code, text)
			return nil
		})
	}
	defer room.Remove(playerId)

	msgCh := make(chan *game.Message)
	go func() {
		for {
			message, err := player.Receive()
			if err != nil {
				log.Println("[ERROR] WSPlayer couldn't receive a message")
				close(msgCh)
				return
			}

			msgCh <- message
		}
	}()
	for {
		err := game.Gameplay(room, playerId, msgCh)
		if err != nil {
			log.Println(err)
			break
		}

		if room.TryRestart() {
			room.Reset()
			room.Start(mydb, rng, qtd)

			message, err := game.NewRestartMsg()
			if err != nil {
				log.Fatalln("Failed to create restart message")
			}

			room.Broadcast(nil, message)
		}
	}
	log.Println("Killing a handle connection goroutine")
}

func handleCheckCanConnectClosure(rooms map[string]*game.Room) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request at /validate")

		if r.Method == http.MethodGet {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			roomId := r.URL.Query().Get("roomId")
			room, ok := rooms[roomId]
			if !ok {
				http.Error(w, "This room doesn't exists", http.StatusNotFound)
				return
			}

			playerId := r.URL.Query().Get("playerId")

			if !room.PlayerExists(playerId) {
				if room.Full() {
					http.Error(w, "The room is full", http.StatusConflict)
					return
				}

				if room.GetStatus() != game.Waiting {
					log.Println("The room isn't at lobby")
					http.Error(w, "The room isn't at lobby", http.StatusConflict)
					return
				}
			} else {
				player := room.GetPlayer(playerId)
				if player.GetState().Connected {
					http.Error(w, "The player is already connected", http.StatusConflict)
				}
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func wsHandlerClosure(rooms map[string]*game.Room, mydb *sql.DB, rng *rand.Rand, qtd int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request at /ws")

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		roomId := r.Form.Get("roomId")
		room, ok := rooms[roomId]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		playerId := r.Form.Get("playerId")

		var player game.Player
		if !room.PlayerExists(playerId) {
			if room.GetStatus() != game.Waiting {
				log.Println("The room isn't at lobby")
				return
			}

			used := make([]string, 0)
			for _, p := range room.Players {
				used = append(used, p.GetName())
			}

			player = &game.WSPlayer{
				Name:  randomColor(used),
				Ready: false,
				Mu:    &sync.Mutex{},
			}

			playerId = game.GeneratePlayerUUID()
			ok = room.Add(playerId, player)
			if !ok {
				log.Println("Couldn't add the player")
				return
			}
			log.Printf("Created player with id %s\n", playerId)
		} else {
			log.Printf("Player %s returned to the game\n", playerId)
			player = room.GetPlayer(playerId)
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading: ", err)
			return
		}

		wsPlayer, ok := player.(*game.WSPlayer)
		if ok {
			log.Println("Player is WsPlayer")
			wsPlayer.Connect(conn)
		}

		go handleConnection(room, playerId, mydb, rng, qtd, func() {
			if room.Empty() {
				log.Printf("Killing room %s\n", roomId)

				delete(rooms, roomId)
			}
		})
	}
}

func handleListRoomsClosure(rooms map[string]*game.Room) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			responses := make([]RoomResponseDTO, 0)
			for id := range rooms {
				responses = append(responses, RoomResponseDTO{
					Id: id,
				})
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")

			body := ListRoomsResponseDTO{
				Rooms: responses,
			}

			err := json.NewEncoder(w).Encode(&body)
			if err != nil {
				log.Println("Error at list rooms")
				return
			}

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func exitCallback() func() error {
	return func() error {
		os.Exit(0)
		return nil
	}
}

func seedCallback() func() error {
	return func() error {
		mydb, err := db.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		return db.Seed(mydb)
	}
}

func sendMessageClosure(host *game.HostPlayer) func(*game.Message) error {
	return func(message *game.Message) error {
		return host.Send(message)
	}
}

func quitClosure() func() error {
	return func() error {
		return nil
	}
}

func startServerCallback() func() error {
	return func() error {
		mydb, err := db.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		rooms := make(map[string]*game.Room)

		config := game.ReadConfig()
		seed := rand.NewSource(42)
		rng := rand.New(seed)

		log.Println("Server start")

		http.HandleFunc("/list", handleListRoomsClosure(rooms))
		http.HandleFunc("/create", handleCreateRoomClosure(rooms, config, mydb, rng, 5))
		http.HandleFunc("/ws", wsHandlerClosure(rooms, mydb, rng, 5))
		http.HandleFunc("/validate", handleCheckCanConnectClosure(rooms))

		fmt.Println("Server listening at port 8080")
		fmt.Println("Routes")
		fmt.Println("/list")
		fmt.Println("/create")
		fmt.Println("/ws")
		fmt.Println("/validate")
		return http.ListenAndServe(":8080", nil)
	}
}

func startHostCallback() func() error {
	return func() error {
		mydb, err := db.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		rooms := make(map[string]*game.Room)
		id := game.GenerateRoomUUID()

		config := game.ReadConfig()
		rooms[id] = game.NewRoom(config)
		room := rooms[id]

		seed := rand.NewSource(42)
		rng := rand.New(seed)

		room.Start(mydb, rng, 5)

		log.Println("Server start")

		http.HandleFunc("/ws", wsHandlerClosure(rooms, mydb, rng, 5))

		hostPlayer := &game.HostPlayer{
			Name:    "Host",
			Channel: make(chan []byte),
			Ready:   false,
		}
		playerId := game.GeneratePlayerUUID()
		room.Add(playerId, hostPlayer)

		myModel := tui.InitModel(sendMessageClosure(hostPlayer), quitClosure())
		p := tea.NewProgram(myModel)

		msgCh := make(chan *game.Message)
		go func() {
			for {
				response, err := hostPlayer.Receive()
				if err != nil {
					log.Println("[ERROR] Host couldn't receive the message")
				}
				log.Printf("Host received: %s", response.Cmd)

				switch response.Cmd {
				case "start":
					msg := game.StartMsg{}
					err = json.Unmarshal(response.Payload, &msg)
					if err != nil {
						p.Send(tui.NotifyMsg{Text: err.Error()})
						continue
					}

					p.Send(tui.StartMsg{Msg: msg})
				case "lobby":
					msg := game.LobbyMsg{}
					err = json.Unmarshal(response.Payload, &msg)
					if err != nil {
						p.Send(tui.NotifyMsg{Text: err.Error()})
						continue
					}

					p.Send(tui.LobbyMsg{Msg: msg})
				case "notify":
					msg := game.NotifyMsg{}
					err = json.Unmarshal(response.Payload, &msg)
					log.Printf("Notification content: %s", msg.Text)
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
				case "playerStatus":
					continue
				case "guessResponse":
					continue
				default:
					msgCh <- response
				}
			}
		}()

		go game.Gameplay(room, playerId, msgCh)

		go func() {
			//fmt.Println("Listening on port 8080")
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
		}

		os.Exit(0)

		return nil
	}
}

func handleCreateRoomClosure(
	rooms map[string]*game.Room,
	config *game.RoomConfig,
	mydb *sql.DB,
	rng *rand.Rand,
	qtd int,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request at /create")
		if r.Method == http.MethodPost {
			// maximo de rooms de uma vez
			if len(rooms) >= 5 {
				http.Error(w, "Maximum number of rooms reached", http.StatusConflict)
				return
			}

			id := game.GenerateRoomUUID()

			room := game.NewRoom(config)
			room.Start(mydb, rng, qtd)

			rooms[id] = room

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")

			response := CreateRoomResponseDTO{
				Id: id,
			}

			err := json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Println("Failed to encode create room response")
			}

			log.Printf("Room created with id: %s\n", id)

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	fmt.Println("Starting")

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	seedCallbackFn := seedCallback()
	startHostCallbackFn := startHostCallback()
	startServerCallbackFn := startServerCallback()
	exitCallbackFn := exitCallback()

	cmds := map[string]struct {
		Name     string
		Callback func() error
	}{
		"seed": {
			Name:     "seed",
			Callback: seedCallbackFn,
		},
		"host": {
			Name:     "host",
			Callback: startHostCallbackFn,
		},
		"server": {
			Name:     "server",
			Callback: startServerCallbackFn,
		},
		"exit": {
			Name:     "exit",
			Callback: exitCallbackFn,
		},
	}

	for {
		var input string

		fmt.Print("> ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println(err)
		}

		cmd, ok := cmds[input]
		if !ok {
			fmt.Println("Invalid command")
			continue
		}

		err = cmd.Callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}
