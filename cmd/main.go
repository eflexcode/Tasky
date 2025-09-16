package main

import (
	"cmd/data"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func main() {

	var taskyValue string = "help"
	var addValue string = "help"

	flag.StringVar(&taskyValue, "tasky", "", "command -tasky help for info")
	flag.StringVar(&addValue, "add", "", "test commcnd")

	flag.Parse()

	// fmt.Print(addValue)

	db, err := sql.Open("sqlite3", "./tasky.db")

	defer db.Close()

	if err != nil {
		fmt.Println("Error", "Failed to load/connect sqlite db")
		return
	}

	dataRepo := data.DataRepository{
		DB: db,
	}

	switch {

	case taskyValue == "help":

		fmt.Println("help:", "prints out all commands and their use case")
		fmt.Println("add:", "create new task")
		fmt.Println("delete:", "delete existing task with all sub tasks")
		fmt.Println("list:", "get all task with all sub tasks")
		fmt.Println("get:", "get task by id")
		fmt.Println("subTask add:", "create subtask by task id")
		fmt.Println("subTask get:", "get subtask by id")
		fmt.Println("subTask delete:", "delete subtask by id")
		fmt.Println("subTask update:", "update task by id")

	case addValue != "add":

		valueKeyMap, err := ParseKeyValue(addValue)

		if err != nil {
			fmt.Println(err)
			return
		}

		title := valueKeyMap["title"]
		description := valueKeyMap["description"]

		if title == "" || description == "" {
			fmt.Println("key title or description is missing")
		}

		dataRepo.InsertTask(title,description)

	case taskyValue == "update":
	case taskyValue == "delete":
	case taskyValue == "list":
	case taskyValue == "get":
	case taskyValue == "subTask add":
	case taskyValue == "subTask get":
	case taskyValue == "subTask delete":
	case taskyValue == "subTask update":

	}
}

func ParseKeyValue(input string) (map[string]string, error) {

	result := make(map[string]string)

	splitKeyValues := strings.Split(input, ";")

	for _, keyValue := range splitKeyValues {

		isKeyValueEmpty := strings.TrimSpace(keyValue)

		if isKeyValueEmpty == "" {
			continue
		}

		keyAndValue := strings.SplitN(keyValue, "=", 2)

		if len(keyAndValue) != 2 {
			return nil, errors.New("invalid key and value pair")
		}

		key := keyAndValue[0]
		value := keyAndValue[1]

		if key == "" {
			return nil, errors.New("invalid key is empty")
		}

		result[key] = value

	}

	return result, nil
}
