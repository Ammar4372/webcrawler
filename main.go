package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

type page struct {
	link    string
	visited int
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no arguments provided")
		os.Exit(1)
	}

	if len(args) < 3 {
		fmt.Println("too few arguments provided")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseUrl := args[0]

	baseUrl, err := url.Parse(rawBaseUrl)
	if err != nil {
		fmt.Println("invalid Url given")
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("second argument must be a integer")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("third argument must be a integer")
		os.Exit(1)
	}
	cfg := config{
		baseUrl:            baseUrl,
		pages:              map[string]int{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           int(maxPages),
	}

	fmt.Printf("Starting crawl of: %s \n", rawBaseUrl)
	cfg.wg.Add(1)
	go cfg.crawlPage(cfg.baseUrl.String())
	cfg.wg.Wait()
	printReport(cfg.pages, rawBaseUrl)
}
func getHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 400 {
		return "", fmt.Errorf("given Url resulted in error code: %d", res.StatusCode)
	}
	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("given Url did not return HTML")
	}
	html, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(html), nil
}

func (cfg *config) crawlPage(rawCurrentUrl string) {
	defer func() {
		cfg.wg.Done()
		<-cfg.concurrencyControl
	}()
	cfg.concurrencyControl <- struct{}{}
	if cfg.getPagesLen() >= cfg.maxPages {
		return
	}
	currentUrl, err := url.Parse(rawCurrentUrl)

	if err != nil {
		fmt.Print("invalid url")
		return
	}
	if cfg.baseUrl.Hostname() != currentUrl.Hostname() {
		return
	}
	norm_url, _ := Normalize_url(rawCurrentUrl)
	if cfg.addPageVisit(norm_url) {
		return
	}
	html, _ := getHTML(rawCurrentUrl)

	fmt.Printf("Crawling URL: %s \n", norm_url)

	urls, _ := getURLsFromHTML(html, cfg.baseUrl.String())
	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
func (cfg *config) getPagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}

func (cfg *config) addPageVisit(normalizedUrl string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, isFirst = cfg.pages[normalizedUrl]; isFirst {
		cfg.pages[normalizedUrl]++
		return
	}
	cfg.pages[normalizedUrl] = 1
	return
}

func printReport(pages map[string]int, baseUrl string) {
	fmt.Printf("==========================\n REPORT for %s \n==========================\n", baseUrl)
	sorted := sortPages(pages)
	for _, v := range sorted {
		fmt.Printf("Found %d internal links to %s\n", v.visited, v.link)
	}
}

func sortPages(pages map[string]int) []page {
	var sortedPages []page
	for k, v := range pages {
		entry := page{
			link:    k,
			visited: v,
		}
		sortedPages = append(sortedPages, entry)
	}
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].visited == sortedPages[j].visited {
			return sortedPages[i].link < sortedPages[j].link
		}
		return sortedPages[i].visited > sortedPages[j].visited
	})
	return sortedPages
}
