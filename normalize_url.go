package main

import (
	"net/url"
	"strings"
)

func normalizeURL(url_string string) (string, error) {
	purl, err := url.Parse(url_string)
	parsed_url_string := purl.Host + strings.TrimRight(purl.Path, "/")
	return parsed_url_string, err
}

func urlDomainsEqual(url_string1 string, url_string2 string) bool {
	// Returns if the urls are on the same domain.
	purl1, err := url.Parse(url_string1)
	if err != nil {
		return false
	}
	purl2, err := url.Parse(url_string2)
	if err != nil {
		return false
	}
	return (purl1.Host == purl2.Host)
}
