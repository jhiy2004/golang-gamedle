package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
)

func handlePlayerReady(room *Room, player Player, msgCh chan *Message) bool {
	return handlePlayerWaitInput(room, player, msgCh, Playing)
}

func handlePlayerRetry(room *Room, player Player, msgCh chan *Message) bool {
	return handlePlayerWaitInput(room, player, msgCh, Waiting)
}

func handlePlayerWaitInput(room *Room, player Player, msgCh chan *Message, nextStatus RoomState) bool {
	select {
	case <-room.WaitReady():
		room.Mu.Lock()
		room.Status = nextStatus
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
			return false
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
				return true
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
				return true
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

	return true
}

func Gameplay(room *Room, player Player, msgCh chan *Message, mydb *sql.DB, rng *rand.Rand, qtd int) {
	for {
		err := GameLobby(room, player, msgCh)
		if err != nil {
			return
		}

		err = GameQuestions(room, player, msgCh)
		if err != nil {
			return
		}

		// TODO: Estudar uma ideia melhor de reiniciar o jogo
		err = GameEnd(room, player, msgCh, mydb, rng, qtd)
		if err != nil {
			return
		}
	}
}

func GameLobby(room *Room, player Player, msgCh chan *Message) error {
	log.Println("Game Lobby")

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

		ok := handlePlayerReady(room, player, msgCh)
		if !ok {
			return errors.New("Error at the lobby")
		}
	}

	return nil
}

func GameQuestions(room *Room, player Player, msgCh chan *Message) error {
	log.Println("Game Questions")

	questionCnt := 0
	qtdQuestions := len(room.QuestionsOrder)

	message, err := NewPlayerStatusMsg(player.GetName(), questionCnt)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	room.Broadcast(nil, message)

	for questionCnt < qtdQuestions {
		currQuestion := room.QuestionsOrder[questionCnt]
		playersNames := make([]string, 0)
		for p, active := range room.Players {
			if active {
				playersNames = append(playersNames, p.GetName())
			}
		}

		message, err := NewStateMsg(
			room.Questions[currQuestion].Question,
			player.GetName(),
			"",
			RoomStateToString(room.Status),
			playersNames,
		)
		if err != nil {
			log.Fatal(err)
		}
		err = player.Send(message)
		if err != nil {
			fmt.Println(err)
		}

		select {
		case <-room.WaitIsEnded():
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				log.Println("Client exit the game early")
				return errors.New("Error at gameplay")
			}

			if msg.Cmd == "guess" {
				log.Println("Guess...")
				guessMsg := GuessMsg{}
				err = json.Unmarshal(msg.Payload, &guessMsg)
				if err != nil {
					log.Fatal(err)
				}

				if room.ValidateAnswer(currQuestion, guessMsg.Answer) {
					// TODO: Remove the notification logic
					message, err := NewNotifyMsg("You're right!")
					if err != nil {
						log.Fatal(err)
					}
					err = player.Send(message)

					message, err = NewGuessResponseMsg(true, "You're right")
					if err != nil {
						log.Fatal(err)
					}

					if err != nil {
						log.Fatal(err)
					}

					// TODO: Remove the notification logic
					message, err = NewNotifyMsg(fmt.Sprintf("User %s passed the %d question", player.GetName(), questionCnt+1))
					if err != nil {
						log.Fatal(err)
					}

					room.Broadcast(
						player,
						message,
					)

					message, err = NewPlayerStatusMsg(player.GetName(), questionCnt+1)
					if err != nil {
						log.Fatal(err)
					}

					room.Broadcast(
						nil,
						message,
					)

					questionCnt++
				} else {
					// TODO: Remove the notification logic
					message, err := NewNotifyMsg("You're wrong haha!")
					if err != nil {
						log.Fatal(err)
					}

					message, err = NewGuessResponseMsg(false, "You're wrong haha!")
					if err != nil {
						log.Fatal(err)
					}

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

	return nil
}

func GameEnd(room *Room, player Player, msgCh chan *Message, mydb *sql.DB, rng *rand.Rand, qtd int) error {
	log.Println("Game End")

	ok := room.EndGame(player)
	if ok {
		playersNames := make([]string, 0)
		for p, active := range room.Players {
			if active {
				playersNames = append(playersNames, p.GetName())
			}
		}
		message, err := NewStateMsg(
			"Game ended!!!",
			player.GetName(),
			player.GetName(),
			RoomStateToString(End),
			playersNames,
		)
		err = player.Send(message)
		if err != nil {
			log.Println("[ERROR] Failed to receive end game state")
			return err
		}
		room.Broadcast(player, message)

		message, err = NewNotifyMsg("You've won the game")
		if err != nil {
			return err
		}

		err = player.Send(message)
		if err != nil {
			log.Println("[ERROR] Failed to win the game")
			return err
		}

		message, err = NewNotifyMsg(fmt.Sprintf("Player %s won the game", player.GetName()))
		if err != nil {
			return err
		}

		room.Broadcast(player, message)

		// O jogador que ganhou a partida se torna responsavel
		// por inicializar a goroutine que aguarda todos clicarem
		// em retry para reiniciar o jogo
		go func() {
			<-room.WaitReady()
			room.Reset()
			room.Start(mydb, rng, qtd)
		}()
	}

	// TODO: Atualmente estou aproveitando as variaveis e logica para ready para implementar o retry
	// Talvez mudar no futuro
	if player.IsReady() {
		player.ToggleReady()
	}

	// Retry
	for room.GetStatus() == End {
		log.Println("Iter")
		ok := handlePlayerRetry(room, player, msgCh)
		if !ok {
			return errors.New("Error at end game")
		}
	}

	return nil
}
