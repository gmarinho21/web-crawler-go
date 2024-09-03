package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(URL string) (string, error) {

	u, err := url.Parse(URL)

	if err != nil {
		fmt.Println(err)
		return "", errors.New("Could not parse URL")
	}

	hostname := u.Host
	path := strings.TrimSuffix(u.Path, "/")

	normalizedURL := hostname + path
	return normalizedURL, err
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var URLs []string

	doc, err := html.Parse(strings.NewReader(htmlBody))

	if err != nil {
		return nil, errors.New("Could not parse html body.")
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		var fullURL string
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					URL, _ := url.Parse(a.Val)
					fullURL = URL.String()

					if URL.Host == "" {
						fullURL = rawBaseURL + URL.Path
					}

					URLs = append(URLs, fullURL)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return URLs, nil
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)

	if err != nil {
		return "", errors.New("could not finish the get request")
	}

	if resp.StatusCode >= 400 {
		return "", errors.New("the request returned with an error")
	}

	contentType := resp.Header["Content-Type"][0]

	if !strings.Contains(contentType, "text/html") {
		fmt.Println(contentType)
		return "", errors.New("the request returned with a not supported content-type, please ensure it return html")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed reading the body of the response")
	}
	resp.Body.Close()

	return string(body), nil
}
