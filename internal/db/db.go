
package db



import (
	"os"
	"log"
	// "fmt"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)




type CrudOperations struct {
	db *sql.DB

}

func Connection() (*sql.DB) {
	var db *sql.DB

	DatabaseConfig := mysql.Config{
		User:     os.Getenv("DBUSER"),
		Passwd:   os.Getenv("DBPASS"),
		Net:      "tcp",
		Addr:     "127.0.0.1:3306",
		DBName:   "tasks",
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

func Connect() CrudOperations{
	operations := CrudOperations{db: Connection()}
	return operations
}