package sitemap

import (
	"bytes"
	"encoding/xml"
	"gophercises/sitemap/link"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func Generate(url string, maxDepth int) (string, error) {

	pages := bfs(url, maxDepth)
	toXml := urlset{
		Xmlns: xmlns,
	}

	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	buf := bytes.NewBufferString(xml.Header)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")
	if err := enc.Encode(toXml); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	// queries
	var q map[string]struct{}
	// next queries, init'd to the original url
	nq := map[string]struct{}{
		urlStr: {},
	}

	for i := 0; i <= maxDepth; i++ {
		// set queries to the next queries, set next queries to a new map
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}

		for url := range q {
			// skip if we've already seen the url
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}
			// add the links from the next page to the next queries map
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}

	return ret
}

func get(urlStr string) []string {
	log.Println("Requesting:", urlStr)
	// HTTP request
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	// get the base domain
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	// use the `links` package to get the links, then filters to the ones with the base domain
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	// get links
	links, _ := link.Parse(r)
	var ret []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			// local link, need to append base url
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			// already has base url
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
