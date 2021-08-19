package db

import (
	"Education/internal/app/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type PostRepository struct {
	dbc *DBController
}

func (repos *PostRepository) FindIfExists(postToFind *models.Post) (int, bool) {
	query := `SELECT id FROM posts WHERE user_id=? AND title=? AND body=?`
	repos.dbc.ctx, repos.dbc.cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer repos.dbc.cancelFunc()
	stmt, err := repos.dbc.db.PrepareContext(repos.dbc.ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var idToReturn int
	err = stmt.QueryRowContext(repos.dbc.ctx, postToFind.UserId, postToFind.Title, postToFind.Body).Scan(&idToReturn)
	if err != nil {
		return 0, false
	}
	return idToReturn, true
}

func (repos *PostRepository) Update(postToUpdate *models.Post, fields map[string]interface{}) error {
	repos.dbc.ctx, repos.dbc.cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer repos.dbc.cancelFunc()
	id, isFind := repos.FindIfExists(postToUpdate)
	if isFind {
		for key := range fields {
			switch key {
			case "user_id":
				postToUpdate.UserId = fields[key].(int)
			case "title":
				postToUpdate.Title = fields[key].(string)
			case "body":
				postToUpdate.Body = fields[key].(string)
			}
		}
	}
	query := "UPDATE posts SET user_id=?,title=?,body=? WHERE id=?"
	stmt, err := repos.dbc.db.PrepareContext(repos.dbc.ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(repos.dbc.ctx, postToUpdate.UserId, postToUpdate.Title, postToUpdate.Body, id)
	if err != nil {
		return err
	}
	return nil
}

func (repos *PostRepository) Insert(postToInsert *models.Post) error {
	query := "INSERT INTO posts (user_id,title,body) VALUES (?,?,?)"
	if id, isFind := repos.FindIfExists(postToInsert); isFind {
		fmt.Println(isFind, id)
		return errors.New("duplicated record")
	}
	repos.dbc.ctx, repos.dbc.cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer repos.dbc.cancelFunc()
	stmt, err := repos.dbc.db.PrepareContext(repos.dbc.ctx, query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(repos.dbc.ctx, postToInsert.UserId, postToInsert.Title, postToInsert.Body)
	if err != nil {
		return err
	}
	return nil
}
