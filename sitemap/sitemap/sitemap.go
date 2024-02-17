package sitemap

import (
	"bytes"
	"encoding/xml"
	"gophercises/sitemap/link"
	"io"
	"log"
	"net/http"
	"strings"
)

func Generate(initial string) (string, error) {
	var urlStack []string
	urlStack = append(urlStack, initial)

	after, _ := strings.CutPrefix(initial, "https://")
	domain := strings.Split(after, "/")[0]

	sitemap := make(map[string]uint8)

	for len(urlStack) > 0 {
		target := pop(&urlStack)

		log.Println("Requesting:", target)
		body, err := get(target)
		if err != nil {
			return "", err
		}

		links, err := links(body)
		if err != nil {
			return "", err
		}

		urlStack = append(urlStack, urls(domain, &sitemap, links)...)
	}

	sitemapXML, err := genSitemapXML(sitemap)
	if err != nil {
		return "", err
	}

	return sitemapXML, nil
}

func validate(domain string, url string) bool {

	// this should catch `mail:` and `file://` urls
	if strings.Contains(strings.Split(url, ".")[0], ":") && !strings.Contains(url, "http") {
		return false
	}

	return strings.Contains(url, domain)
}

func get(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func links(r io.Reader) ([]link.Link, error) {
	links, err := link.Parse(r)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func urls(domain string, sitemap *map[string]uint8, links []link.Link) []string {
	urlStack := make([]string, 0)

	for _, link := range links {

		cleaned := clean(domain, link.Href)
		if _, ok := (*sitemap)[cleaned]; ok {
			continue
		}

		if ok := validate(domain, cleaned); !ok {
			continue
		}

		// the list doesn't exist on the map but is valid
		(*sitemap)[cleaned] = 0
		urlStack = append(urlStack, cleaned)
	}

	return urlStack
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url
}

func genSitemapXML(urls map[string]uint8) (string, error) {
	set := &UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  make([]Url, 0),
	}

	// urls is a map so we're getting the key
	for url := range urls {
		set.Urls = append(set.Urls, Url{Loc: url})
	}

	buf := bytes.NewBufferString("")
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")
	if err := enc.Encode(set); err != nil {
		return "", err
	}

	xmlString := strings.Join([]string{xml.Header, buf.String()}, "")

	return xmlString, nil
}

func pop[T any](slice *[]T) T {
	x := (*slice)[0]
	*slice = (*slice)[1:]
	return x
}

func clean(domain string, link string) string {
	if strings.Contains(link, ":") {
		return link
	} else {
		return strings.Join([]string{"https://", domain, link}, "")
	}
}
