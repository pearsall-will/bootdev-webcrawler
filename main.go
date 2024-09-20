package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %v\n", args[0])
	// html, err := getHTML(args[0])
	// if err != nil {
	// 	fmt.Println(err)}
	var waitGroup sync.WaitGroup

	cfg := config{
		pages:   map[string]int{},
		baseURL: args[0],
		mu:      &sync.Mutex{},
		wg:      &waitGroup,
		cc:      make(chan struct{}, 1),
	}
	cfg.crawlPage(args[0])
	cfg.wg.Wait()
	// fmt.Println(html)
}
