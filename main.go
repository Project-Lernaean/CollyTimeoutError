package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"CollyError/shared"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Open URL List
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	var URLS []string

	for scanner.Scan() {
		URLS = append(URLS, scanner.Text())
	}

	// Create Colly Collector
	c := colly.NewCollector(
		colly.MaxBodySize(10e9),
		colly.DetectCharset(),
		colly.Async(true),
		colly.IgnoreRobotsTxt(),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 250})
	c.WithTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		DialContext: (&net.Dialer{
			Timeout:  5 * time.Second,
			Deadline: time.Now().Add(20 * time.Second),
		}).DialContext,
	})

	// Create Results Channel
	pageInfoResults := make(chan shared.PageInfo, 10000)
	pageHTMLResults := make(chan shared.HtmlInfo, 10000)
	pagecount := 0

	// Define Process for Responses
	c.OnResponse(func(r *colly.Response) {
		pagecount++
		fmt.Println(pagecount)
		u, err := url.Parse(fmt.Sprintf("%s", r.Request.URL))
		if err != nil {
			// log.Println("OnResponse:", err)
		}

		Port, err := strconv.Atoi(u.Port())
		if err != nil {
			// log.Println("OnResponse:", err)
		}

		pi := shared.PageInfo{
			Port:        Port,
			StatusCode:  r.StatusCode,
			Length:      len(r.Body),
			ServerType:  r.Headers.Get("server"),
			ContentType: r.Headers.Get("content-type"),
			Headers:     make(map[string]string),
		}

		// Collect Headers
		for name, values := range *r.Headers {
			for _, value := range values {
				pi.Headers[name] = value
				break
			}
		}

		pageInfoResults <- pi
	})

	// Define Process of HTML Page
	c.OnHTML("title", func(e *colly.HTMLElement) {
		_key := GetDBKey(e.Request.URL)
		hi := shared.HtmlInfo{Key: _key, Title: e.Text}
		pageHTMLResults <- hi
	})

	// Define Process of Error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(pagecount)
		pagecount++
	})

	// Spawn Jobs and Wait
	log.Println("Starting Job")
	log.Println("Attempting", len(URLS), "URLS")
	for _, x := range URLS {
		c.Visit(fmt.Sprintf("http://%s", x))
	}
	c.Wait()

	// Parse Results
	log.Println("Finished Job")
	close(pageHTMLResults)
	close(pageInfoResults)
	cache := make(map[string]shared.PageInfo)

	// Assmeble Responses and HTML
	for pi := range pageInfoResults {
		cache[pi.Key] = pi
	}
	for hi := range pageHTMLResults {
		pi := cache[hi.Key]
		pi.Title = hi.Title
		cache[hi.Key] = pi
	}
	fmt.Println(cache)

}

func GetDBKey(u *url.URL) string {
	return fmt.Sprintf("%s_%s", u.Hostname(), u.Port())
}
