package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(inputURL string, inputBody string) ([]string, error) {
	tokens := html.NewTokenizer(strings.NewReader(inputBody))
	var links []string
	for {
		tokenType := tokens.Next()
		if tokenType == html.ErrorToken {
			err := tokens.Err()
			if err == io.EOF {
				break
			}
		}
		tokn := tokens.Token()
		if tokn.Data == "a" {
			for _, atr := range tokn.Attr {
				if atr.Key != "href" {
					continue
				}
				if !strings.Contains(atr.Val, "http") {
					links = append(links, inputURL+atr.Val)
				} else {
					links = append(links, atr.Val)
				}

			}
		}
	}
	return links, nil
}
