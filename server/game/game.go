package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhiy2004/golang-gamedle/server/db"
	"log"
	"math/rand"
)

func GameStart(room *Room, mydb *sql.DB, rng *rand.Rand, qtd int) error {
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
		room.Questions[q.Id] = q
	}
	room.QuestionsOrder = append(room.QuestionsOrder, choices...)

	return nil
}

func Gameplay(room *Room, player Player, msgCh chan *Message) {
	GameLobby(room, player, msgCh)
	GameQuestions(room, player, msgCh)
	GameEnd(room, player)
}

func GameLobby(room *Room, player Player, msgCh chan *Message) {
	message, err := NewStartMsg(room.MinPlayers, room.MaxPlayers, player.GetName())
	if err != nil {
		log.Fatal(err)
	}
	err = player.Send(message)
	if err != nil {
		log.Println(err)
	}
	player.Send(message)

	message, err = NewLobbyMsg(room.CurrPlayers, room.ReadyPlayers)
	if err != nil {
		log.Fatal(err)
	}
	err = player.Send(message)
	if err != nil {
		log.Println(err)
	}
	room.Broadcast(nil, message)

	for room.GetStatus() == Waiting {
		room.WaitMinReached()

		select {
		case <-room.WaitReady():
			room.Mu.Lock()
			room.Status = Playing
			room.Mu.Unlock()

			message, err := NewNotifyMsg("Game is about to start")
			if err != nil {
				log.Fatal(err)
			}
			room.Broadcast(player, message)
		case msg, ok := <-msgCh:
			log.Println(msg)
			if !ok {
				log.Println("OMG")
				return
			}

			if msg.Cmd == "ready" {
				log.Println("Ready message")
				changed := room.PlayerReady(player)

				if !changed {
					message, err := NewNotifyMsg("You already was ready")
					if err != nil {
						log.Fatal(err)
					}
					player.Send(message)
					continue
				}

				message, err := NewNotifyMsg(fmt.Sprintf("Player %s is ready", player.GetName()))
				if err != nil {
					log.Fatal(err)
				}
				room.Broadcast(player, message)
			} else if msg.Cmd == "cancel" {
				log.Println("Cancel message")

				changed := room.PlayerCancel(player)

				if !changed {
					message, err := NewNotifyMsg("You already was not ready")
					if err != nil {
						log.Fatal(err)
					}
					player.Send(message)
					continue
				}

				message, err := NewNotifyMsg(fmt.Sprintf("Player %s cancelled ready operation", player.GetName()))
				if err != nil {
					log.Fatal(err)
				}
				room.Broadcast(player, message)
			}

			message, err := NewLobbyMsg(room.CurrPlayers, room.ReadyPlayers)
			if err != nil {
				log.Fatal(err)
			}
			err = player.Send(message)
			if err != nil {
				log.Println(err)
			}
			room.Broadcast(nil, message)
		}
	}
}

func GameQuestions(room *Room, player Player, msgCh chan *Message) {
	questionCnt := 0
	qtdQuestions := len(room.QuestionsOrder)
	for questionCnt < qtdQuestions {
		currQuestion := room.QuestionsOrder[questionCnt]
		message, err := NewStateMsg(room.Questions[currQuestion].Question, player.GetName(), "", RoomStateToString(room.Status))
		if err != nil {
			log.Fatal(err)
		}
		err = player.Send(message)
		if err != nil {
			fmt.Println(err)
		}

		msg, ok := <-msgCh
		if !ok {
			log.Println("Client exit the game early")
			return
		}

		if msg.Cmd == "guess" {
			log.Println("Guess...")
			guessMsg := GuessMsg{}
			err = json.Unmarshal(msg.Payload, &guessMsg)
			if err != nil {
				log.Fatal(err)
			}

			if room.ValidateAnswer(currQuestion, guessMsg.Answer) {
				message, err := NewNotifyMsg("You're right!")
				if err != nil {
					log.Fatal(err)
				}

				err = player.Send(message)
				if err != nil {
					log.Fatal(err)
				}

				message, err = NewNotifyMsg(fmt.Sprintf("User %s passed the %d question", player.GetName(), questionCnt+1))
				if err != nil {
					log.Fatal(err)
				}

				room.Broadcast(
					player,
					message,
				)
				questionCnt++
			} else {
				message, err := NewNotifyMsg("You're wrong haha!")
				if err != nil {
					log.Fatal(err)
				}

				err = player.Send(message)
				if err != nil {
					log.Fatal(err)
				}

			}
		}

	}
}

func GameEnd(room *Room, player Player) {
	room.Winner = player
	room.Status = End

	message, err := NewStateMsg("Game ended!!!", player.GetName(), player.GetName(), RoomStateToString(End))
	err = player.Send(message)
	if err != nil {
		log.Println("[ERROR] Failed to receive end game state")
		return
	}
	room.Broadcast(player, message)

	message, err = NewNotifyMsg("You've won the game")
	if err != nil {
		log.Fatal(err)
	}

	err = player.Send(message)
	if err != nil {
		log.Println("[ERROR] Failed to win the game")
		return
	}

	message, err = NewNotifyMsg(fmt.Sprintf("Player %s won the game", player.GetName()))
	if err != nil {
		log.Fatal(err)
	}

	room.Broadcast(player, message)
}
