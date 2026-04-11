package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type QuestionEntity struct {
	Id int
	Question string
}

type AnswerEntity struct {
	Id int
	QuestionId int
	Answer string
}

type QuestionAnswersDTO struct {
	Id int
	Question string
	Answers []string
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./clidle.db")	
	if err != nil {
		return nil, err
	}
	
	sqlStmt := `
		create table if not exists question(
			id integer not null primary key,
			question text not null
		)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `
		create table if not exists answer(
			id integer not null primary key,
			question_id integer references question(id),
			answer text not null
		)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetQuestions(db *sql.DB, questionIds []int) ([]QuestionAnswersDTO, error) {
	if questionIds == nil {
		return nil, errors.New("nil value passed as questionIds")
	}

	res := make([]QuestionAnswersDTO, 0)

	qtdIds := len(questionIds)
	var queryIds strings.Builder
	for i, id := range questionIds {
		if i == (qtdIds - 1) {
			fmt.Fprintf(&queryIds, "%d", id)
		} else {
			fmt.Fprintf(&queryIds, "%d,", id)
		}
	}

	stmt := fmt.Sprintf("SELECT id, question FROM question WHERE id IN (%s)", queryIds.String())
	questionsRows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	rowsCount := 0
	for {
		var (
			id int
			q string
		)

		if !questionsRows.Next() {
			break
		}

		if err := questionsRows.Scan(&id, &q); err != nil {
			return nil, err
		}

		res = append(res, QuestionAnswersDTO{
			Id: id,
			Question: q,
		})

		rowsCount++
	}

	if len(questionIds) != rowsCount {
		return nil, errors.New("At least one question_id doesn't exists in the db")
	}

	for i := range res {
		stmt := fmt.Sprintf("SELECT answer FROM answer WHERE question_id = %d", res[i].Id)
		rows, err := db.Query(stmt)
		if err != nil {
			return nil, err
		}

		answers := make([]string, 0)
		for {
			var ans string

			if !rows.Next() {
				break
			}

			err := rows.Scan(&ans)
			if err != nil {
				return nil, err
			}

			answers = append(answers, ans)

		}

		res[i].Answers = answers
	}

	return res, nil
}

func InsertQuestion(db *sql.DB, question string, answers []string) (*QuestionAnswersDTO, error) {
	var id int


	stmt := "INSERT INTO question (question) VALUES (?) RETURNING id"
	row := db.QueryRow(stmt, question)

	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	for _, ans := range answers {
		stmt = "INSERT INTO answer (question_id, answer) VALUES (?, ?)"
		_, err := db.Exec(stmt, id, ans)
		if err != nil {
			return nil, err
		}
	}

	res := &QuestionAnswersDTO{
		Id: id,
		Question: question,
		Answers: answers,
	}

	return res, nil
}

func GetQuestionsIds(db *sql.DB) ([]int, error) {
	ids := make([]int, 0)

	stmt := "SELECT id FROM question"
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	for {
		var id int

		if !rows.Next() {
			break
		}

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)

	}

	return ids, nil
}


func Seed(db *sql.DB) error {
	batch := []struct{
		q string
		ans []string
	} {
		{
			q: "A plumber who saves a princess",
			ans: []string{"Super Mario", "Mario"},
		},
		{
			q: "A monster-catching RPG",
			ans: []string{"Pokemon", "Pokémon"},
		},
		{
			q: "A block-stacking puzzle game",
			ans: []string{"Tetris"},
		},
		{
			q: "A battle royale with building mechanics",
			ans: []string{"Fortnite"},
		},
		{
			q: "A sandbox game about surviving zombies",
			ans: []string{"Project Zomboid"},
		},
		{
			q: "A sci-fi FPS about fighting aliens",
			ans: []string{"Halo"},
		},
		{
			q: "A post-apocalyptic RPG in a nuclear wasteland",
			ans: []string{"Fallout"},
		},
		{
			q: "A fantasy RPG with dragons and an open world",
			ans: []string{"Skyrim", "The Elder Scrolls V: Skyrim"},
		},
		{
			q: "A stealth game about a genetically engineered assassin",
			ans: []string{"Hitman"},
		},
		{
			q: "A portal-based puzzle game with sarcastic AI",
			ans: []string{"Portal", "Portal 2"},
		},
		{
			q: "A MOBA with champions and summoners",
			ans: []string{"League of Legends", "LoL"},
		},
		{
			q: "A survival game about exploring an alien ocean",
			ans: []string{"Subnautica"},
		},
		{
			q: "A roguelike dungeon crawler about escaping the underworld",
			ans: []string{"Hades"},
		},
		{
			q: "A rhythm game where you slash blocks with lightsabers",
			ans: []string{"Beat Saber"},
		},
	}

	for _, entry := range batch {
		_, err := InsertQuestion(db, entry.q, entry.ans)
		if err != nil {
			return err
		}
	}

	return nil
}
