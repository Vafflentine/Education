package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

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

// WriteToConsole Output post to console
func WriteToConsole(response []byte) {
	if _, err := fmt.Println(string(response)); err != nil {
		log.Fatal(err)
	}
}

// WriteToFile Create file and write post to this file
func WriteToFile(response []byte, countFile int) {
	fileName := "storage/posts/" + strconv.Itoa(countFile) + ".txt"
	if err := ioutil.WriteFile(fileName, response, 0600); err != nil {
		log.Fatal(err)
	}
}

func main() {
	getPostsAsync(100, "console")
}
