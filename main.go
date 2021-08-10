package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

type Post struct {
	Id      int    `json:"id"`
	User_id int    `json:"userId"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

type Comment struct {
	Id      int    `json:"id"`
	Post_id int    `json:"postId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Body    string `json:"body"`
}

// Getting all posts sync
func getPosts(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body)) // ???
	return body
}

// Getting all posts async
func getPostsAsync(count int, destination string) {
	wg := new(sync.WaitGroup)
	for i := 1; i <= count; i++ {
		wg.Add(1)
		i := i
		go func() {
			response := getPosts("https://jsonplaceholder.typicode.com/posts/" + strconv.Itoa(i))
			switch destination {
			case "console":
				WriteToConsole(response)
			case "file":
				WriteToFile(response, i)
			}
			defer wg.Done()
		}()
	}

	wg.Wait()
}

// Output post to console
func WriteToConsole(response []byte) {
	if _, err := fmt.Println(string(response)); err != nil {
		log.Fatal(err)
	}
}

// Create file and write post to this file
func WriteToFile(response []byte, countFile int) {
	fileName := "storage/posts/" + strconv.Itoa(countFile) + ".txt"
	if err := ioutil.WriteFile(fileName, response, 0600); err != nil {
		log.Fatal(err)
	}
}

func getComments(post []Post, comment []Comment) {
	wg := new(sync.WaitGroup)
	for _, i := range post {
		wg.Add(1)
		i := i
		if err := Insert(i, "posts"); err != nil {
			log.Fatal(err)
		}
		go func() {
			response := getPosts("https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(i.Id))
			if err := json.Unmarshal(response, &comment); err != nil {
				log.Fatal(err)
			}

			for j := range comment {
				j := j
				go func() {
					if err := Insert(comment[j], "comments"); err != nil {
						log.Fatal(err)
					}
				}()
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func ConnectToDB(connectionString string) error {
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		return pingErr
	}
	fmt.Println("Connected to DB")
	return nil
}

func Insert(data interface{}, tableName string) error {
	switch obj := data.(type) {
	case Post:
		_, err := db.Exec("INSERT INTO "+tableName+" (id,user_id,title,body) VALUES (?,?,?,?)", &obj.Id, &obj.User_id, &obj.Title, &obj.Title)
		if err != nil {
			return err
		}
	case Comment:
		_, err := db.Exec("INSERT INTO "+tableName+" (id,post_id,name,email,body) VALUES (?,?,?,?,?)", &obj.Post_id, &obj.Id, &obj.Name, &obj.Email, &obj.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var post []Post
	var comment []Comment
	if err := ConnectToDB("root:@/education"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	response := getPosts("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err := json.Unmarshal(response, &post); err != nil {
		log.Fatal(err)
	}
	getComments(post[:], comment[:])
}
