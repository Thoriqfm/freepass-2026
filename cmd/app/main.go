package main

import (
	"freepass-2026/internal/handler/rest"
	"freepass-2026/internal/repository"
	"freepass-2026/internal/service"
	"freepass-2026/pkg/bcrypt"
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

	repo := repository.NewRepository(db)
	bcrypt := bcrypt.Init()
	svc := service.NewService(repo, bcrypt)
	r := rest.NewRest(svc)
	r.MountEndPoint()
	r.Run()
}
