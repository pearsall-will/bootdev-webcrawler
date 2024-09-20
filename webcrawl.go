package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type config struct {
	pages   map[string]int
	baseURL string
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	cc      chan struct{}
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 400 {
		return "", fmt.Errorf("HTTP Error \nStatus Code:%v", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("invalid content-type found %v", resp.Header.Get("content-type"))
	}
	html_text, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(html_text), nil
}

func (conf *config) safePageUpdate(nurl string) bool {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	if val, ok := conf.pages[nurl]; !ok {
		conf.pages[nurl] = 1
	} else {
		conf.pages[nurl] = val + 1
		fmt.Println("Already crawled.")
		return true
	}
	return false
}

func (conf *config) crawlPage(rawCurrentURL string) {
	conf.wg.Add(1)
	conf.cc <- struct{}{}
	defer conf.wg.Done()
	defer func() {
		<-conf.cc
	}()

	// defer _ := <- conf.cc
	if !urlDomainsEqual(conf.baseURL, rawCurrentURL) {
		fmt.Printf("%v not in same domain as %v. Going back.\n", rawCurrentURL, conf.baseURL)
		return
	}

	nurl, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Crawling: %v\n", rawCurrentURL)
	done := conf.safePageUpdate(nurl)
	if done {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}

	links, err := getURLsFromHTML(conf.baseURL, html)
	if err != nil {
		fmt.Println(err)
	}

	for _, link := range links {
		fmt.Println("got to for loop") // If I remove this line.. it doesn't crawl the page, but it does with it in.
		go conf.crawlPage(link)
	}
}
