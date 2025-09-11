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
	SubTask     []SubTask
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

	_,err := data.DB.Exec(query,&task.Title,&task.Description,updatedAt,&task.ID)

	if err != nil {
		return err
	}
	return nil
}
func (data *DataRepository) GetAllTask(task *Task) (error, []Task ){


    query := `SELECT * FROM tasks WHERE id = ?`

	updatedAt := time.Now()

	_,err := data.DB.Exec(query,&task.Title,&task.Description,updatedAt,&task.ID)

	if err != nil {
		return err
	}
	return nil
}

