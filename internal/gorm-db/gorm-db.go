package gormdb

import (
	"crud/internal/db"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrmCrudOperations struct {
	db *gorm.DB
}

func (c *OrmCrudOperations) ListAll() ([]db.Task, error) {

	var tasks []db.Task

	result := c.db.Find(&tasks)

	if err := result.Error; err != nil {

		return nil, fmt.Errorf("ListAll: %v", err)
	}

	return tasks, nil

}

func (c *OrmCrudOperations) GetTaskById(Id int) (db.Task, error) {
	var task db.Task

	result := c.db.First(&task, Id)

	if err := result.Error; err != nil {
		return task, fmt.Errorf("GetTaskById: %v", err)
	}

	return task, nil
}

func (c *OrmCrudOperations) GetTaskByCompletion(completed bool) ([]db.Task, error) {

	var task []db.Task
	result := c.db.Where("completed = ?", completed).Find(&task)

	if err := result.Error; err != nil {

		return task, fmt.Errorf("GetTaskByCompletion: %v", err)
	}

	return nil, fmt.Errorf("error")

}

func (c *OrmCrudOperations) UpdateTask(taskId int64, task db.Task) (int64, error) {

	var modifiedTask db.Task

	if err := c.db.Model(&modifiedTask).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where("id = ?", taskId).Updates(db.Task{Name: task.Name, Completed: task.Completed}).Error; err != nil {

		return 0, fmt.Errorf("UpdateTask: %v", err)
	}

	return int64(modifiedTask.Id), nil
}

func (c *OrmCrudOperations) CreateTask(task db.Task) (int64, error) {
	result := c.db.Create(&task)

	if err := result; err != nil {

		return 0, fmt.Errorf("CreateTask: %v", err)
	}

	return result.RowsAffected, nil

}

func (c *OrmCrudOperations) DeleteTask(taskId int) (int64, error) {

	err := c.db.Delete(&db.Task{}, taskId).Error

	if err != nil {

		return 0, fmt.Errorf("DeleteTask: %v", err)
	}

	return 0, fmt.Errorf("error")

}

func Connection() *gorm.DB {
	dsn := "root:senhafacil@tcp(127.0.0.1:3306)/tasks?charset=utf8mb4&parseTime=True&loc=Local"
	database, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return database
}

func Connect() *OrmCrudOperations {
	operations := &OrmCrudOperations{Connection()}
	return operations
}
