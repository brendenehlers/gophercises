package main

import (
	"flag"
	"fmt"
	"gophercises/sitemap/sitemap"
)

func main() {
	siteFlag := flag.String("site", "https://google.com", "Site to generate a sitemap for")
	depthFlag := flag.Int("depth", 10, "The maximum depth to search to")
	flag.Parse()

	xml, err := sitemap.Generate(*siteFlag, *depthFlag)
	if err != nil {
		panic(err)
	}

	fmt.Println(xml)
}
