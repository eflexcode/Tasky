package main

import (
	"cmd/data"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func main() {

	db, err := sql.Open("sqlite3", "./tasky.db?_busy_timeout=10000&_journal_mode=WAL")

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

	querySub := `CREATE TABLE IF NOT EXISTs subtasks(id INTEGER PRIMARY KEY AUTOINCREMENT, TaskID INTEGER, title TEXT,description TEXT,done BOOLEAN ,created_at DATETIME ,updated_at DATETIME,FOREIGN KEY (TaskID) REFERENCES tasks(id) ON DELETE CASCADE )`

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
	var deleteValue int = -8
	var getValue int = -1
	var toggleValue string = "toggle"
	var updateValue string = "update"

	flag.BoolVar(&taskyValue, "help", false, "command -tasky help for info")
	flag.StringVar(&addValue, "add", "add", "create new task")
	flag.BoolVar(&listValue, "list", false, "show all task and subtask")
	flag.IntVar(&deleteValue, "delete", -8, "delete with id")
	flag.IntVar(&getValue, "get", -1, "get with id")
	flag.StringVar(&toggleValue, "toggle", "toggle", "toggle if task is done")
	flag.StringVar(&updateValue, "update", "update", "update task title and description with id")

	flag.Parse()

	switch {

	case taskyValue:

		fmt.Println("Syntax(put in double quote):", "id=1;done=yes\n")

		fmt.Println("-help:", "prints out all commands and their use case\n")
		fmt.Println("-add:", "create new task(title and description)\n")
		fmt.Println("-delete:", "delete existing task with all sub tasks\n")
		fmt.Println("-list:", "get all task with all sub tasks\n")
		fmt.Println("-update:", "update task by id\n")
		fmt.Println("-toggle:", "update task done by id\n")
		fmt.Println("-get:", "get task by id\n")
		// fmt.Println("subTask add:", "create subtask by task id")
		// fmt.Println("subTask get:", "get subtask by id")
		// fmt.Println("subTask delete:", "delete subtask by id")
		// fmt.Println("subTask update:", "update task by id")

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

	case updateValue != "update":

		valueKeyMap, err := ParseKeyValue(updateValue)

		if err != nil {
			fmt.Println(err)
			return
		}

		id := valueKeyMap["id"]
		title := valueKeyMap["title"]
		description := valueKeyMap["description"]

		idInt, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println("Error", "Id is not a number")
			return
		}

		t, err := dataRepo.GetTaskById(int64(idInt))

		if err != nil {
			fmt.Println("Error", "No task found with id "+id)
			return
		}

		if title != "" {
			t.Title = title
		}

		if description != "" {
			t.Description = description
		}

		fmt.Println(t.Title + " " + t.Description)

		err = dataRepo.UpdateTask(t.Title,t.Description,int64(idInt))

		if err != nil {
			fmt.Println(err)
			fmt.Println("Error", "Failed to update task "+id)
			return
		}

		fmt.Println("Success", "Task updated succesfully")

	case deleteValue != -8:
		err := dataRepo.DeleteTask(int64(deleteValue))
		if err != nil {
			fmt.Println("Error", "Failed to delete task")
			return
		}

		fmt.Println("Success", "Task with id: "+strconv.Itoa(deleteValue)+" deleted succesfullly")

	case listValue:

		tasks, err := dataRepo.GetAllTask()

		if err != nil {
			fmt.Print("Error", "Failed to load task from db")
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 2, '*', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintln(w, "ID\tTitle\tDescription\tDone\tCreatedAt\tUpdatedAt\t")

		for _, t := range tasks {

			var doneS string

			if t.Done {
				doneS = "YES"
			} else {
				doneS = "NO"
			}

			// var subTaskDoneCount int

			// for _, d := range t.SubTasks {

			// 	if d.Done {
			// 		subTaskDoneCount++
			// 	}

			// }

			// subTaskInfo := "T: " + strconv.Itoa(len(t.SubTasks)) + " D: " + strconv.Itoa(subTaskDoneCount)

			tString := strconv.Itoa(t.ID) + "\t" + t.Title + "\t" + t.Description +  "\t" + doneS + "\t" + t.CreatedAt.Format("02/01/2006 15:04:05") + "\t" + t.UpdatedAt.Format("02/01/2006 15:04:05") + "\t"

			fmt.Fprintln(w, tString)

		}

		fmt.Fprintln(w)
		w.Flush()

	case getValue != -1:

		t, err := dataRepo.GetTaskById(int64(getValue))

		if err != nil {
			fmt.Println("Error", "Failed to get task with id "+strconv.Itoa(getValue))
			return
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 2, '*', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintln(w, "ID\tTitle\tDescription\tDone\tCreatedAt\tUpdatedAt\t")
		var doneS string

		if t.Done {
			doneS = "YES"
		} else {
			doneS = "NO"
		}

		// var subTaskDoneCount int

		// for _, d := range t.SubTasks {

		// 	if d.Done {
		// 		subTaskDoneCount++
		// 	}

		// }

		// subTaskInfo := "T: " + strconv.Itoa(len(t.SubTasks)) + " D: " + strconv.Itoa(subTaskDoneCount)

		tString := strconv.Itoa(t.ID) + "\t" + t.Title + "\t" + t.Description  + "\t" + doneS + "\t" + t.CreatedAt.Format("02/01/2006 15:04:05") + "\t" + t.UpdatedAt.Format("02/01/2006 15:04:05") + "\t"

		fmt.Fprintln(w, tString)
		fmt.Fprintln(w)
		w.Flush()
	case toggleValue != "toggle":

		valueKeyMap, err := ParseKeyValue(toggleValue)

		if err != nil {
			fmt.Println(err)
			return
		}

		done := valueKeyMap["done"]
		id := valueKeyMap["id"]

		if done == "" {
			fmt.Println("Error", "key done or id is missing")
			return
		}

		var isDone bool

		if strings.EqualFold(done, "Yes") {
			isDone = true
		} else {
			isDone = false
		}

		v, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println("Error", "id is not a number")
			return
		}

		err = dataRepo.ToggleTask(isDone, int64(v))

		if err != nil {
			fmt.Println("Error", "Faild to update task")
			return
		}
		fmt.Println("Success", "Task updated succesfully")

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
