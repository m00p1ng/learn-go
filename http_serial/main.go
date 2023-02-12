package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func checkUrl(url string, c chan string) {
	resp, err := http.Get(url)

	if err != nil {
		// fmt.Println(err)
		s := fmt.Sprintf("%s is DOWN!\n", url)
		s += fmt.Sprintf("Error: %v\n", err)
		fmt.Println(s)
		c <- url // sending into the channel
	} else {
		s := fmt.Sprintf("%s -> Status Code: %d \n", url, resp.StatusCode)
		s += fmt.Sprintf("%s is UP\n", url)
		fmt.Println(s)
		c <- url
	}
}

func main() {
	urls := []string{
		"https://golang.org",
		"https://www.google.com",
		"https://www.medium.com",
	}

	c := make(chan string)

	for _, url := range urls {
		go checkUrl(url, c)
	}

	fmt.Println("No. of Goroutines: ", runtime.NumGoroutine())

	for {
		go checkUrl(<-c, c)
		fmt.Println(strings.Repeat("#", 30))
		time.Sleep(time.Second * 2)
	}
}
