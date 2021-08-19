package main

import (
	"Education/internal/app/apiserver"
	"Education/internal/app/models"
	"log"
)

func main() {
	config := apiserver.NewConfig()
	server := apiserver.New(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
	err := server.DBController.Post().Insert(&models.Post{UserId: 1, Title: "test", Body: "body"})
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()
}
