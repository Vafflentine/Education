package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func main() {
	fmt.Println(string(getPosts("https://jsonplaceholder.typicode.com/posts")))
}
