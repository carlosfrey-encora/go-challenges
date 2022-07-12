package api

import (
	crud_via_gorm "crud/internal/crud-via-gorm"
	crud_via_vanilla "crud/internal/crud-via-vanilla"
	"crud/internal/grpc/pb"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CrudProvider interface {
	ListAll() ([]*pb.Task, error)
	GetTaskById(Id int) (*pb.Task, error)
	GetTaskByCompletion(completed bool) ([]*pb.Task, error)
	UpdateTask(taskId int64, task *pb.Task) (int64, error)
	CreateTask(task *pb.Task) (int64, error)
	DeleteTask(taskId int) (int64, error)
}

type ApiService struct {
	dbService CrudProvider
}

func (a *ApiService) GetAll(w http.ResponseWriter, r *http.Request) {

	tasks, err := a.dbService.ListAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(tasks)
	}

}

func (a *ApiService) GetById(w http.ResponseWriter, r *http.Request, id int) {

	task, err := a.dbService.GetTaskById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(task)
	}

}

func (a *ApiService) GetByCompletion(w http.ResponseWriter, r *http.Request, completed bool) {

	tasks, err := a.dbService.GetTaskByCompletion(completed)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(tasks)
	}

}

func (a *ApiService) PostTask(w http.ResponseWriter, r *http.Request) {
	var task pb.Task
	json.NewDecoder(r.Body).Decode(&task)

	id, err := a.dbService.CreateTask(&task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(id)
	}

}

func (a *ApiService) PutTask(w http.ResponseWriter, r *http.Request, taskId int64) {

	var task pb.Task

	json.NewDecoder(r.Body).Decode(&task)

	_, err := a.dbService.UpdateTask(taskId, &task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Task sucessfully changed!")
	}

}

func (a *ApiService) DeleteTask(w http.ResponseWriter, r *http.Request, taskId int) {

	id, err := a.dbService.DeleteTask(taskId)

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(id)
	}

}

func AllowAnyOrigin(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (a *ApiService) MainHandler(w http.ResponseWriter, r *http.Request) {

	AllowAnyOrigin(&w)

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
		apiService = ApiService{crud_via_vanilla.Connect()}
	} else if implementation == "gorm" {
		apiService = ApiService{crud_via_gorm.Connect()}
	}

	http.HandleFunc("/tasks/", apiService.MainHandler)
	http.ListenAndServe(":8000", http.DefaultServeMux)

}
