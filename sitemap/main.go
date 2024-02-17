package main

import (
	"flag"
	"fmt"
	"gophercises/sitemap/sitemap"
	"strings"
)

func main() {
	siteFlag := flag.String("site", "https://google.com", "Site to generate a sitemap for")
	flag.Parse()

	var url string
	if strings.HasPrefix(*siteFlag, "https://") {
		url = *siteFlag
	} else {
		url = strings.Join([]string{"https://", *siteFlag}, "")
	}

	xml, err := sitemap.Generate(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(xml)
}
