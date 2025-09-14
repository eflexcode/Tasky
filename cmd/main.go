package main

import (
	"flag"
	"fmt"
)

func main() {

	var taskyValue string = "help"

	flag.StringVar(&taskyValue, "tasky", "help", "command -tasky help for info")

	flag.Parse()

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

	case taskyValue == "add":
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

func f(p *string) {

}
