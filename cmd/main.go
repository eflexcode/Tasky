package main

import (
	"cmd/data"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"strings"
)

func main() {

	db, err := sql.Open("sqlite3", "./tasky.db")

	defer db.Close()

	if err != nil {
		fmt.Println("Error", "Failed to load/connect sqlite db")
		return
	}

	query := `CREATE TABLE IF NOT EXISTs tasks(id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT,description TEXT,done BOOLEAN ,created_at DATETIME ,updated_at DATETIME )`

	_, err = db.Exec(query)

	if err != nil {
		fmt.Println("Error", "Failed to load/connect sqlite db table")
		return
	}

	querySub := `CREATE TABLE IF NOT EXISTs subtasks(id INTEGER PRIMARY KEY AUTOINCREMENT, TaskID INTEGER, title TEXT,description TEXT,done BOOLEAN ,created_at DATETIME ,updated_at DATETIME )`

	_, err = db.Exec(querySub)

	if err != nil {
		fmt.Println("Error", "Failed to load/connect sqlite db table")
		return
	}

	dataRepo := data.DataRepository{
		DB: db,
	}

	var taskyValue bool
	var addValue string = "add"
	var listValue bool = false

	flag.BoolVar(&taskyValue, "help", false, "command -tasky help for info")
	flag.StringVar(&addValue, "add", "add", "create new task")
	flag.BoolVar(&listValue, "list", false, "show all task and subtask")
	flag.String("ok", "skks", "ksksk")

	flag.Parse()

	switch {

	case taskyValue:

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
			return
		}

		err = dataRepo.InsertTask(title, description)

		if err != nil {
			fmt.Println("Error", "error while trying to isert to db")
		} else {
			fmt.Println("Success", "Task with title "+title+" and description "+description+" added")
		}

	// case taskyValue == "update":
	// case taskyValue == "delete":
	case listValue:

		tasks,err := dataRepo.GetAllTask()

		if err != nil {
			fmt.Print("Error","Failed to load task from db")
		}

		

	// case taskyValue == "get":
	// case taskyValue == "subTask add":
	// case taskyValue == "subTask get":
	// case taskyValue == "subTask delete":
	// case taskyValue == "subTask update":

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
