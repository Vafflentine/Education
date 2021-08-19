package db

import (
	"Education/internal/app/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type DBController struct {
	config            *DBConfig
	db                *sql.DB
	ctx               context.Context
	cancelFunc        context.CancelFunc
	postRepository    *PostRepository
	commentRepository *CommentRepository
}

func New(config *DBConfig) *DBController {
	return &DBController{
		config: config,
	}
}

func (dbc *DBController) NewConnection() error {
	dbc.config.Dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbc.config.DBUser,
		dbc.config.DBPassword,
		dbc.config.Address,
		dbc.config.DBPort,
		dbc.config.DBName)
	db, err := sql.Open("mysql", dbc.config.Dns)
	if err != nil {
		log.Fatalf("Error %s when opening DB\n", err)
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	dbc.ctx, dbc.cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer dbc.cancelFunc()
	err = db.PingContext(dbc.ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return err
	}

	fmt.Printf("Connected to DB %s successfully\n", dbc.config.Dns)

	dbc.db = db
	if err := dbc.migrate(models.Post{}, models.Comment{}); err != nil {
		return err
	}
	return nil
}

func (dbc *DBController) CloseConnection() error {
	if err := dbc.db.Close(); err != nil {
		return err
	}
	return nil
}

func (dbc *DBController) Post() *PostRepository {
	if dbc.postRepository != nil {
		return dbc.postRepository
	}
	dbc.postRepository = &PostRepository{
		dbc: dbc,
	}
	return dbc.postRepository
}

func (dbc *DBController) Comment() *CommentRepository {
	if dbc.commentRepository != nil {
		return dbc.commentRepository
	}
	dbc.commentRepository = &CommentRepository{
		dbc: dbc,
	}
	return dbc.commentRepository
}

func (dbc *DBController) migrate(modelsToMigrate ...interface{}) error {
	dbc.ctx, dbc.cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer dbc.cancelFunc()
	for _, value := range modelsToMigrate {
		switch value.(type) {
		case models.Post:
			query := `CREATE TABLE IF NOT EXISTS posts (
			id int(11) PRIMARY KEY AUTO_INCREMENT,
			user_id int(11) NOT NULL,
			title text NOT NULL,
			body text NOT NULL,
			created_at datetime(3) DEFAULT CURRENT_TIMESTAMP,
			updated_at datetime(3) DEFAULT CURRENT_TIMESTAMP) 
			ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

			_, err := dbc.db.ExecContext(dbc.ctx, query)
			if err != nil {
				return err
			}
		case models.Comment:
			query := `CREATE TABLE IF NOT EXISTS comments (
			post_id int(11) NOT NULL,
			id int(11) PRIMARY KEY AUTO_INCREMENT,
			name varchar(60) NOT NULL,
			email varchar(40) NOT NULL,
			body text NOT NULL,
			created_at datetime(3) DEFAULT CURRENT_TIMESTAMP,
			updated_at datetime(3) DEFAULT CURRENT_TIMESTAMP)
			ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

			_, err := dbc.db.ExecContext(dbc.ctx, query)
			if err != nil {
				return err
			}
		default:
			return errors.New("can't migrate any table")
		}
	}
	return nil
}
