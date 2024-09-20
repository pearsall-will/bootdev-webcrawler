package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type config struct {
	pages    map[string]int
	baseURL  string
	mu       *sync.Mutex
	wg       *sync.WaitGroup
	cc       chan struct{}
	maxPages int
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
	// Guard rail to dead stop on max depth reached.
	if len(conf.pages) >= conf.maxPages {
		// fmt.Println("Max Depth Reached")
		return true
	}

	val, ok := conf.pages[nurl]

	if !ok {
		conf.pages[nurl] = 1
		return false
	}

	conf.pages[nurl] = val + 1
	// fmt.Println("Already crawled.")
	return true

}

func (conf *config) maxDepthReached() bool {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	return len(conf.pages) >= conf.maxPages
}

func (conf *config) crawlPage(rawCurrentURL string) {
	defer conf.wg.Done()
	if conf.maxDepthReached() {
		return
	}

	defer func() {
		<-conf.cc
	}()
	conf.cc <- struct{}{}

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
		conf.wg.Add(1)
		go conf.crawlPage(link)
	}
}

func sortMap(pages map[string]int) {
	pairs := make([][2]interface{}, 0, len(pages))
	for k, v := range pages {
		pairs = append(pairs, [2]interface{}{k, v})
	}

	// Sort slice based on values
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][1].(int) < pairs[j][1].(int)
	})

	// Extract sorted keys
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p[0].(string)
	}

	// Print sorted map
	for _, k := range keys {
		fmt.Printf("Found %d: internal links to https://%s\n", pages[k], k)
	}
}

func (conf *config) printReport() {
	fmt.Println("=============================")
	fmt.Printf("  REPORT for %s\n", conf.baseURL)
	fmt.Println("=============================")
	sortMap(conf.pages)
}
