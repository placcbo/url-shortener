package main

import "fmt"

type Link struct {
	Code   string
	URL    string
	Clicks int
}

func main() {

	link := Link{
		Code:   "abc123",
		URL:    "http://go.dev",
		Clicks: 0,
	}

	fmt.Println(link)
}