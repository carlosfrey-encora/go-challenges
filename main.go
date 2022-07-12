package main

import (
	"crud/internal/api"
	"os"
)

func main() {

	os.Setenv("DB_IMPL", "orm")
	api.SetupApi()
}
