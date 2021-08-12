package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Post struct {
	gorm.Model
	Id      int    `json:"id"`
	User_id int    `json:"userId"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

type Comment struct {
	gorm.Model
	Id      int    `json:"id"`
	Post_id int    `json:"postId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Body    string `json:"body"`
}

var (
	db *gorm.DB
)

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

func getComments(post []Post, comment []Comment) {
	wg := new(sync.WaitGroup)
	for _, i := range post {
		wg.Add(1)
		go func(i Post) {
			db.Create(&Post{Id: i.Id, User_id: i.User_id, Title: i.Title, Body: i.Body})
			/*if err := Insert(i, "posts"); err != nil {
				log.Fatal(err)
			}*/
			response := getPosts("https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(i.Id))
			if err := json.Unmarshal(response, &comment); err != nil {
				log.Fatal(err)
			}
			for _, j := range comment {
				go func(j Comment) {
					db.Create(&Comment{Id: j.Id, Post_id: j.Post_id, Name: j.Name, Email: j.Email, Body: j.Body})
				}(j)
			}
			defer wg.Done()
		}(i)
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

func ConnectToDB(connectionString string) error {
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	fmt.Println("Connected")
	return nil
}

func Migrate(post *Post, comment *Comment) error {
	if postErr := db.AutoMigrate(&post); postErr != nil {
		return postErr
	}
	if commentErr := db.AutoMigrate(&comment); commentErr != nil {
		return commentErr
	}
	return nil
}

func main() {
	var post []Post
	var comment []Comment
	if err := ConnectToDB("root:@tcp(127.0.0.1:3306)/education?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		log.Fatal(err)
	}
	// Close
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Migrate schemas
	if err := Migrate(&Post{}, &Comment{}); err != nil {
		log.Fatal(err)
	}

	response := getPosts("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err := json.Unmarshal(response, &post); err != nil {
		log.Fatal(err)
	}
	getComments(post[:], comment[:])

}
