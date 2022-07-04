package main



import (
	"fmt"
	"crud/internal/db"
)



func main() {
	dbConnection := db.Connect()
	allTasks, _ := dbConnection.GetTaskByCompletion(false)

	fmt.Println(allTasks)
}