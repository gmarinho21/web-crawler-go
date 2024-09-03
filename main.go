package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	inputArgs := os.Args[1:]

	if len(inputArgs) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(inputArgs) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	BASE_URL := inputArgs[0]

	fmt.Println("starting crawl of: ", BASE_URL)
	// html, _ := getHTML(BASE_URL)

	// fmt.Println(html)
	pages := make(map[string]int)
	crawlPage(BASE_URL, BASE_URL, pages)
	fmt.Println(pages)
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	fmt.Printf("--------- New Request %s  -------\n", rawCurrentURL)
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if baseURL.Host != currentURL.Host {
		fmt.Printf("URL Not in same domain as base. \nBase: %s\nCurrent: %s\n", baseURL.Host, currentURL.Host)
		return
	}

	normCurrentUrl, err := normalizeURL(rawCurrentURL)

	if err != nil {
		fmt.Println(err)
		return
	}

	if pages[normCurrentUrl] > 0 {
		pages[normCurrentUrl] += 1
		fmt.Println("Page already crawled")
		return
	}
	pages[normCurrentUrl] = 1

	requestedURL := "https://" + normCurrentUrl
	fmt.Println(requestedURL)
	body, err := getHTML(requestedURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	fetchedURLs, err := getURLsFromHTML(body, rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, url := range fetchedURLs {
		crawlPage(rawBaseURL, url, pages)
	}

}
