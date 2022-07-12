package main



import (
	"fmt"
	"crud/internal/db"
)



func main() {
	dbConnection := db.Connect()
	fmt.Println(dbConnection)
}