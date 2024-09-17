package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlbody, baseUrl string) ([]string, error) {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL %w", err)
	}
	reader := strings.NewReader(htmlbody)
	node, err := html.Parse(reader)
	var links []string
	if err != nil {
		return nil, fmt.Errorf("error parsing html %w", err)
	}
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, v := range n.Attr {
				if v.Key == "href" {
					URL, err := url.Parse(v.Val)
					if err != nil {
						continue
					}
					resolvedURL := baseURL.ResolveReference(URL)
					links = append(links, resolvedURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return links, nil
}
