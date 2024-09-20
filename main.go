package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	url := args[0]
	parallelism, err := strconv.Atoi(args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(maxPages)
	fmt.Println(parallelism)
	fmt.Printf("starting crawl of: %v\n", args[0])
	var wg sync.WaitGroup

	cfg := config{
		pages:    map[string]int{},
		baseURL:  args[0],
		mu:       &sync.Mutex{},
		wg:       &wg,
		cc:       make(chan struct{}, parallelism),
		maxPages: maxPages,
	}
	cfg.wg.Add(1)
	cfg.crawlPage(url)
	cfg.wg.Wait()
    cfg.printReport()
}
