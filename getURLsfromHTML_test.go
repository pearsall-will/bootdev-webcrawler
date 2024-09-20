package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `<!DOCTYPE html>
				<html>
				<body>

				<h2>HTML Links</h2>
				<p>HTML links are defined with the a tag:</p>

				<a href="https://www.w3schools.com">This is a link</a>

				</body>
				</html>`,
			expected: []string{"https://www.w3schools.com"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://go.dev/tour/moretypes/15",
			inputBody: `<!doctype html>
<html lang="en" ng-app="tour" data-theme="auto">

<head>
    <meta charset="utf-8">
    <title>A Tour of Go</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="mobile-web-app-capable" content="yes">
    <script>
      (function() {
        const theme = document.cookie.match(/prefers-color-scheme=(light|dark|auto)/)?.[1]
        if (theme) {
          document.querySelector('html').setAttribute('data-theme', theme);
        }
      }())
    </script>
    <link rel="stylesheet" href="/tour/static/css/app.css" />
    <link rel="stylesheet" href="/tour/static/lib/codemirror/lib/codemirror.css">
    <link href='https://fonts.googleapis.com/css?family=Inconsolata' rel='stylesheet' type='text/css'>
    <link rel="icon" href="/images/favicon-gopher.png" sizes="any">
    <link rel="apple-touch-icon" href="/images/favicon-gopher-plain.png"/>
    <link rel="icon" href="/images/favicon-gopher.svg" type="image/svg+xml">
</head>

<body>
    <div class="bar top-bar">
        <div class="left">
        <a href="/"><img src="/images/go-logo-white.svg" class="gopherlogo"></a>
        <a class="logo" href="/tour/list">A Tour of Go</a>
        </div>
        <div class="right">
            <button class="header-toggleTheme js-toggleTheme" aria-label="Toggle theme">
              <img
                data-value="auto"
                class="go-Icon go-Icon--inverted"
                height="24"
                width="24"
                src="/images/icons/brightness_6_gm_grey_24dp.svg"
                alt="System theme">
              <img
                data-value="dark"
                class="go-Icon go-Icon--inverted"
                height="24"
                width="24"
                src="/images/icons/brightness_2_gm_grey_24dp.svg"
                alt="Dark theme">
              <img
                data-value="light"
                class="go-Icon go-Icon--inverted"
                height="24"
                width="24"
                src="/images/icons/light_mode_gm_grey_24dp.svg"
                alt="Light theme">
            </button>
            <div table-of-contents-button=".toc"></div>
            <div feedback-button></div>
        </div>
    </div>

    <div table-of-contents></div>

    <div ng-view ng-cloak class="ng-cloak"></div>

    <script src="/tour/script.js"></script>
</body>

</html>`,
			expected: []string{"https://go.dev/tour/moretypes/15/", "https://go.dev/tour/moretypes/15/tour/list"},
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputURL, tc.inputBody)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
