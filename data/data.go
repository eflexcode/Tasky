package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	SubTasks    []SubTask
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubTask struct {
	ID        int
	TaskID    int
	Title     string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DataRepository struct {
	DB *sql.DB
}

func NewDataRepository(DB *sql.DB) *DataRepository {
	return &DataRepository{DB: DB}
}

func (data *DataRepository) InsertTask(task *Task) error {

	query := `INSERT INTO tasks(title,description,created_at,updated_at) VALUES($1,$2,$3,$4)`

	createdAt := time.Now()

	_, err := data.DB.Exec(query, &task.Title, &task.Description, createdAt, createdAt)

	if err != nil {
		return err
	}
	return nil
}

func (data *DataRepository) UpdateTask(task *Task) error {

	query := `UPDATE tasks SET title = ?,description = ?,updated_at = ? WHERE id = ?`

	updatedAt := time.Now()

	_, err := data.DB.Exec(query, &task.Title, &task.Description, updatedAt, &task.ID)

	if err != nil {
		return err
	}

	return nil
}

func (data *DataRepository) GetAllTask() ([]Task, error) {

	var tasks []Task

	query := `SELECT * FROM tasks`

	row, err := data.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {

		item := Task{}
		var subItems []SubTask

		err := row.Scan(&item.ID, &item.Title, &item.Description, &item.UpdatedAt, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		querySubTask := `SELECT * FROM subtasks WHERE TaskID = ?`

		rowSubTask, err := data.DB.Query(querySubTask, item.ID)

		if err != nil {
			return nil, err
		}

		for rowSubTask.Next() {
			subItem := SubTask{}
			err = rowSubTask.Scan(&subItem.ID, &subItem.TaskID, &subItem.Title, &subItem.Done, &subItem.CreatedAt, &subItem.UpdatedAt)

			if err != nil {
				return nil, err
			}

			subItems = append(subItems, subItem)
			item.SubTasks = subItems
		}

		tasks = append(tasks, item)
	}

	return tasks, nil
}

func (data *DataRepository) GetTaskById(id int64) (*Task, error) {

	query := `SELECT * FROM tasks WHERE id = ?`

	result, err := data.DB.Query(query, id)

	if err != nil {
		return nil, err
	}

	task := Task{}

	err = result.Scan(&task.ID, &task.Title, &task.Description, &task.UpdatedAt, &task.CreatedAt)

	if err != nil {
		return nil, err
	}

	querySubTask := `SELECT * FROM subtasks WHERE TaskID = ?`

	rowSubTask, err := data.DB.Query(querySubTask, id)

	if err != nil {
		return nil, err
	}

	var subTasks []SubTask

	for rowSubTask.Next() {

		subTask := SubTask{}

		err = rowSubTask.Scan(&subTask.ID, &subTask.TaskID, &subTask.Title, &subTask.Done, &subTask.CreatedAt, &subTask.UpdatedAt)

		if err != nil {
			return nil, err
		}

		subTasks = append(subTasks, subTask)

	}

	task.SubTasks = subTasks

	return &task, nil

}

func (data *DataRepository) DeleteTask(id int64) error {

	query := `DELETE FROM tasks WHERE id = ?`
	_, err := data.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (data *DataRepository) AddSubTask(taskId int64, title string) (*Task, error) {

	query := `INSERT INTO subtasks(TaskID,title,done,created_at,updated_at) VALUES($1,$2,$3,$4,$5)`
	createdAt := time.Now()

	_, err := data.DB.Exec(query, taskId, title, false, createdAt, createdAt)

	if err != nil {
		return nil, err
	}

	task, err := data.GetTaskById(taskId)

	if err != nil {
		return nil, err
	}

	return task, nil

}

func (data *DataRepository) UpdateSubTask(subTaskId int64) {}

func (data *DataRepository) DeleteSubTask(id int64) error {

	query := `DELETE FROM subtasks WHERE id = ?`
	_, err := data.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
