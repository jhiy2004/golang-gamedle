package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func handlePlayerRetry(room *Room, player Player, msg *Message) {
	if msg.Cmd == "retry" {
		log.Println("Retry Message")
		changed := room.PlayerRetry(player)

		if !changed {
			message, err := NewNotifyMsg("You already was retry")
			if err != nil {
				log.Fatal(err)
			}
			player.Send(message)
		}

		message, err := NewNotifyMsg(fmt.Sprintf("Player %s is retry", player.GetName()))
		if err != nil {
			log.Fatal(err)
		}
		room.Broadcast(player, message)
	} else if msg.Cmd == "cancelRetry" {
		log.Println("Cancel Retry message")

		changed := room.PlayerCancelRetry(player)

		if !changed {
			message, err := NewNotifyMsg("You already was not retry")
			if err != nil {
				log.Fatal(err)
			}
			player.Send(message)
		}

		message, err := NewNotifyMsg(fmt.Sprintf("Player %s cancelled retry operation", player.GetName()))
		if err != nil {
			log.Fatal(err)
		}
		room.Broadcast(player, message)
	}

	message, err := NewPostGameLobbyMsg(room.CurrPlayers, room.RetryPlayers)
	if err != nil {
		log.Fatal(err)
	}
	err = player.Send(message)
	if err != nil {
		log.Println(err)
	}

	log.Println(message)
	room.Broadcast(nil, message)
}

func handlePlayerReady(room *Room, player Player, msg *Message) {
	if msg.Cmd == "ready" {
		log.Println("Ready message")
		changed := room.PlayerReady(player)

		if !changed {
			message, err := NewNotifyMsg("You already was ready")
			if err != nil {
				log.Fatal(err)
			}
			player.Send(message)
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

	log.Println(message)
	room.Broadcast(nil, message)
}

func Gameplay(room *Room, player Player, msgCh chan *Message) error {
	err := GameLobby(room, player, msgCh)
	if err != nil {
		return err
	}

	err = GameQuestions(room, player, msgCh)
	if err != nil {
		return err
	}

	err = GameEnd(room, player, msgCh)
	if err != nil {
		return err
	}

	return nil
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

	for {
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
			return nil

		case msg, ok := <-msgCh:
			if !ok {
				return errors.New("Error at the lobby")
			}
			handlePlayerReady(room, player, msg)
		}
	}
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
		for _, p := range room.Players {
			playersNames = append(playersNames, p.GetName())
		}

		message, err := NewStateMsg(
			room.Questions[currQuestion].Question,
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

func GameEnd(room *Room, player Player, msgCh chan *Message) error {
	log.Println("Game End")

	ok := room.EndGame(player)
	if ok {
		playersNames := make([]string, 0)
		for _, p := range room.Players {
			playersNames = append(playersNames, p.GetName())
		}
		message, err := NewStateMsg(
			"Game ended!!!",
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
	}

	// Retry
	retryCh := room.WaitRetry()
	for {
		log.Println("Iter")

		select {
		case <-retryCh:
			return nil

		case msg, ok := <-msgCh:
			if !ok {
				return errors.New("Error at end game")
			}
			handlePlayerRetry(room, player, msg)
		}
	}
}
