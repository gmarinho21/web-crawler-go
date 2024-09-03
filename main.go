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

	const maxConcurrency = 3
	cfg, err := configure(BASE_URL, maxConcurrency)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	fmt.Println("starting crawl of: ", BASE_URL)

	cfg.wg.Add(1)
	go cfg.crawlPage(BASE_URL)
	cfg.wg.Wait()

	for normalizedURL, count := range cfg.pages {
		fmt.Printf("%d - %s\n", count, normalizedURL)
	}
}

func (cfg *config) crawlPage(rawCurrentURL string) {

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	fmt.Printf("--------- New Request %s  -------\n", rawCurrentURL)
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		fmt.Printf("URL Not in same domain as base. \nBase: %s\nCurrent: %s\n", cfg.baseURL.Host, currentURL.Host)
		return
	}

	normCurrentUrl, err := normalizeURL(rawCurrentURL)

	if err != nil {
		fmt.Println(err)
		return
	}

	isFirstTimeRequesting := cfg.addPageVisit(normCurrentUrl)
	if !isFirstTimeRequesting {
		return
	}

	requestedURL := "https://" + normCurrentUrl
	fmt.Println(requestedURL)
	body, err := getHTML(requestedURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	fetchedURLs, err := getURLsFromHTML(body, cfg.baseURL.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, url := range fetchedURLs {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}

}
