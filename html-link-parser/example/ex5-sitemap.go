package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zhenyang0405/gophercises/html-link-parser/parser"
)

/*
	1. GET the webpage
	2. Parse all the links on the page
	3. Build proper urls with the links
	4. Filter out any links from difference domain
	5. FInd all the pages - Breath First Search
	6. Print our using XML
*/

const xmlNameSpace = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls []loc `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://www.kristenleecalligraphy.com/", "url to get to build sitemap")
	maxDepth := flag.Int("depth", 3, "the number to traverse")
	flag.Parse()

	//fmt.Println(*urlFlag)

	pages := depthFirstSearch(*urlFlag, *maxDepth)
	toXml := urlset {
		Xmlns: xmlNameSpace,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	fmt.Print(xml.Header)
	encode := xml.NewEncoder(os.Stdout)
	encode.Indent("", "  ")
	if err := encode.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

func depthFirstSearch(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})
		if len(queue) == 0 {
			break
		}
		for webUrl, _ := range queue {
			if _, ok := seen[webUrl]; ok {
				continue
			}
			seen[webUrl] = struct{}{}
			for _, link := range getUrl(webUrl) {
				if _, ok := seen[link]; ok {
					nextQueue[link] = struct{}{}
				}
			}
		}
	}

	arrLists := make([]string, 0, len(seen))
	for webUrl, _ := range seen {
		arrLists = append(arrLists, webUrl)
	}
	return arrLists
}

func getUrl(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host: reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var arrLinks []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			arrLinks = append(arrLinks, base + l.Href)
		case strings.HasPrefix(l.Href, "http"):
			arrLinks = append(arrLinks, l.Href)
		}
	}
	return arrLinks
}

func filter(links []string, keepFn func(string) bool) []string {
	var arrLinks []string
	for _, link := range links {
		if keepFn(link) {
			arrLinks = append(arrLinks, link)
		}
	}
	return arrLinks
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}