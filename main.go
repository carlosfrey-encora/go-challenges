package main

import (
	api "crud/internal/rest-api"
	"os"
)

func main() {

	os.Setenv("DB_IMPL", "gorm")
	api.SetupApi()
}
