package orm

// import (
// 	"fmt"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// 	"crud/internal/db"
// )

// func Connect() {

// 	dsn := "root:senhafacil@tcp(127.0.0.1:3306)/tasks?charset=utf8mb4&parseTime=True&loc=Local"
// 	database, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	var tasks []db.Task

// 	result := database.Find(&tasks)

// 	for _, task := range tasks {

// 		fmt.Println(task)
// 	}

// 	fmt.Println(result)
// 	fmt.Println(result.Error)
// 	fmt.Println(result.RowsAffected)

// }
