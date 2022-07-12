package crud_via_vanilla

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	pb "crud/internal/grpc/pb"

	"github.com/go-sql-driver/mysql"
)

type CrudOperations struct {
	db *sql.DB
}

func (c *CrudOperations) ListAll() ([]*pb.Task, error) {

	var tasks []*pb.Task

	rows, err := c.db.Query("SELECT * FROM task")

	if err != nil {
		return nil, fmt.Errorf("ListAll %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		var task pb.Task

		if err := rows.Scan(&task.Id, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("ListAll: %v", err)
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListAll: %v", err)
	}

	return tasks, nil
}

func (c *CrudOperations) GetTaskById(Id int) (*pb.Task, error) {

	var task pb.Task
	row := c.db.QueryRow("SELECT * FROM task WHERE id = ?", Id)

	if err := row.Scan(&task.Id, &task.Name, &task.Completed); err != nil {

		if err == sql.ErrNoRows {
			return &task, fmt.Errorf("GetTaskById %d: no such album", Id)
		}

		return &task, fmt.Errorf("GetTaskById %d: %v", Id, err)
	}

	return &task, nil

}

func (c *CrudOperations) GetTaskByCompletion(completed bool) ([]*pb.Task, error) {

	var tasks []*pb.Task
	rows, err := c.db.Query("SELECT * FROM task WHERE Completed = ?", completed)

	if err != nil {
		fmt.Errorf("GetTaskByCompletion: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task pb.Task

		if err := rows.Scan(&task.Id, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("GetTaskByCompletion: %v", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTaskByCompletion: %v", err)
	}

	return tasks, nil

}

func (c *CrudOperations) UpdateTask(taskId int64, task *pb.Task) (int64, error) {

	result, err := c.db.Exec("UPDATE task SET name = ?, completed = ? WHERE id = ?", task.Name, task.Completed, taskId)

	if err != nil {
		return 0, fmt.Errorf("UpdateTask: %v", err)
	}

	id, err := result.RowsAffected()

	if id == 0 {
		return 0, fmt.Errorf("UpdateTask: task id wasn't found: %d", taskId)
	}

	if err != nil {
		return 0, fmt.Errorf("UpdateTask: %v", err)
	}

	return id, nil
}

func (c *CrudOperations) CreateTask(task *pb.Task) (int64, error) {
	result, err := c.db.Exec("INSERT INTO task (name, completed) values (?, ?)", task.Name, task.Completed)

	if err != nil {

		return 0, fmt.Errorf("CreateTask: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("CreateTask: %v", err)
	}

	return id, nil

}

func (c *CrudOperations) DeleteTask(taskId int) (int64, error) {

	result, err := c.db.Exec("DELETE FROM task WHERE id = ? ", taskId)

	if err != nil {
		return 0, fmt.Errorf("DeleteTask: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("DeleteTask: %v", err)
	}

	return id, nil

}

func Connection() *sql.DB {
	var db *sql.DB

	DatabaseConfig := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "tasks",
	}

	var err error

	db, err = sql.Open("mysql", DatabaseConfig.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}

func Connect() *CrudOperations {
	operations := CrudOperations{db: Connection()}
	return &operations
}
