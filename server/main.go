package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
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

func randomColor() string {
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

	i := rand.Intn(len(colors))

	return colors[i]
}

func handleConnection(room *game.Room, player *game.WSPlayer) {
	conn := player.Conn

	defer room.Remove(player)
	defer conn.Close()

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Oh shit: %d %s\n", code, text)
		return nil
	})

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

	game.Gameplay(room, player, msgCh)
}

func wsHandlerClosure(room *game.Room) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Println("Error upgrading: ", err)
			return
		}

		player := &game.WSPlayer{
			Conn:  conn,
			Name:  randomColor(),
			Ready: false,
			Mu:    &sync.Mutex{},
		}

		ok := room.Add(player)
		if !ok {
			conn.WriteMessage(websocket.TextMessage, []byte("The room is full, sorry!"))
			conn.Close()
		}

		go handleConnection(room, player)
	}
}

func exitCallback() func() error {
	return func() error {
		os.Exit(0)
		return nil
	}
}

func seedCallback(mydb *sql.DB) func() error {
	return func() error {
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

func startServerCallback(room *game.Room, mydb *sql.DB, rng *rand.Rand, config *game.RoomConfig) func() error {
	return func() error {
		err := game.GameStart(room, mydb, rng, config.QuestionsCount)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Server start")

		http.HandleFunc("/ws", wsHandlerClosure(room))

		hostPlayer := &game.HostPlayer{
			Name:    "Host",
			Channel: make(chan []byte),
			Ready:   false,
		}
		room.Add(hostPlayer)

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
				default:
					msgCh <- response
				}
			}
		}()

		go game.Gameplay(room, hostPlayer, msgCh)

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

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	mydb, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	config := game.ReadConfig()
	room := game.NewRoom(config)

	seed := rand.NewSource(42)
	rng := rand.New(seed)

	seedCallbackFn := seedCallback(mydb)
	startServerCallbackFn := startServerCallback(room, mydb, rng, config)

	exitCallbackFn := exitCallback()

	cmds := map[string]struct {
		Name     string
		Callback func() error
	}{
		"seed": {
			Name:     "seed",
			Callback: seedCallbackFn,
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
