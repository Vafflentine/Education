package apiserver

import (
	"Education/internal/app/db"
	"errors"
	"fmt"
)

type APIServer struct {
	config       *ServerConfig
	DBController *db.DBController
}

func New(config *ServerConfig) *APIServer {
	return &APIServer{
		config: config,
	}
}

func (server *APIServer) Start() error {
	if server.config == nil {
		return errors.New("can't find .env file for configuration")
	}
	if err := server.startDB(); err != nil {
		return err
	}
	fmt.Println("Server has been started")
	return nil
}

func (server *APIServer) startDB() error {
	dbConn := db.New(server.config.DBConfig)
	if err := dbConn.NewConnection(); err != nil {
		return err
	}
	server.DBController = dbConn
	return nil
}

func (server *APIServer) Close() {
	if err := server.DBController.CloseConnection(); err != nil {
		fmt.Println("Can not close connection to Db")
	}
	fmt.Println("Connection has been closed")
}
