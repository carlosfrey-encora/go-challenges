package api

import (
	"crud/internal/db"
	gormdb "crud/internal/gorm-db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CrudProvider interface {
	ListAll() ([]db.Task, error)
	GetTaskById(Id int) (db.Task, error)
	GetTaskByCompletion(completed bool) ([]db.Task, error)
	UpdateTask(taskId int64, task db.Task) (int64, error)
	CreateTask(task db.Task) (int64, error)
	DeleteTask(taskId int) (int64, error)
}

type ApiService struct {
	dbService CrudProvider
}

func (a *ApiService) GetAll(w http.ResponseWriter, r *http.Request) {

	tasks, err := a.dbService.ListAll()

	if err != nil {

		fmt.Errorf("GetAll: %v", err)
	}

	json.NewEncoder(w).Encode(tasks)

}

func (a *ApiService) GetById(w http.ResponseWriter, r *http.Request, id int) {

	task, err := a.dbService.GetTaskById(id)

	if err != nil {

		fmt.Errorf("GetbyId: %v", err)
	}

	if task.Id == 0 {
		fmt.Fprint(w, "There's no task with this id")
	} else {
		json.NewEncoder(w).Encode(task)
	}
}

func (a *ApiService) GetByCompletion(w http.ResponseWriter, r *http.Request, completed bool) {

	tasks, err := a.dbService.GetTaskByCompletion(completed)

	if err != nil {

		fmt.Errorf("GetByCompletion: %v", err)
	}

	json.NewEncoder(w).Encode(tasks)

}

func (a *ApiService) PostTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	json.NewDecoder(r.Body).Decode(&task)

	id, err := a.dbService.CreateTask(task)

	if err != nil {

		fmt.Errorf("PostTask: %v", err)
	}

	json.NewEncoder(w).Encode(id)

}

func (a *ApiService) PutTask(w http.ResponseWriter, r *http.Request, taskId int64) {

	var task db.Task

	json.NewDecoder(r.Body).Decode(&task)

	_, err := a.dbService.UpdateTask(taskId, task)

	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprint(w, "Task sucessfully changed!")
	}

}

func (a *ApiService) DeleteTask(w http.ResponseWriter, r *http.Request, taskId int) {

	id, err := a.dbService.DeleteTask(taskId)

	if err != nil {

		json.NewEncoder(w).Encode(&err)
	} else {

		json.NewEncoder(w).Encode(id)
	}

}

func (a *ApiService) MainHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		if id, length := GetRequestId(r.URL.Path); length > 0 {
			a.GetById(w, r, id)
		} else {
			if param := r.URL.Query().Get("completed"); param != "" {
				isCompleted, _ := strconv.ParseBool(param)
				a.GetByCompletion(w, r, isCompleted)

			} else {

				a.GetAll(w, r)
			}

		}

	case http.MethodPost:
		a.PostTask(w, r)

	case http.MethodPut:

		if id, length := GetRequestId(r.URL.Path); length > 0 {
			a.PutTask(w, r, int64(id))
		}

	case http.MethodDelete:

		if id, length := GetRequestId(r.URL.Path); length > 0 {

			a.DeleteTask(w, r, id)
		}
	}

}

func GetRequestId(url string) (int, int) {

	idString := strings.TrimPrefix(url, "/tasks/")
	idLength := len(idString)
	id, _ := strconv.Atoi(idString)

	return id, idLength
}

func SetupApi() {
	var apiService ApiService
	implementation := os.Getenv("DB_IMPL")

	if implementation == "vanilla" {
		apiService = ApiService{db.Connect()}
	} else if implementation == "orm" {
		apiService = ApiService{gormdb.Connect()}
	}

	http.HandleFunc("/tasks/", apiService.MainHandler)
	http.ListenAndServe(":8000", nil)

}
