package main

import (
	"freepass-2026/pkg/config"
	"freepass-2026/pkg/database"
	"log"
)

func main() {
	config.LoadEnvironment()

	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
}
