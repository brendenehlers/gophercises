package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	filepathFlag := flag.String("file", "samples/ex1.html", "file to read")
	flag.Parse()

	file, err := os.Open(*filepathFlag)
	check(err)
	reader := bufio.NewReader(file)

	doc, err := html.Parse(reader)
	check(err)

	links := findLinksInDocument(doc)
	fmt.Println("Links: ", links)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func findLinksInDocument(doc *html.Node) []Link {

	links := make([]Link, 0)

	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			var href string
			for _, attr := range node.Attr {
				fmt.Println(attr.Key)
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			// TODO get the text using a DFS on the children of this node
			text := "dummy text"

			if href != "" {
				links = append(links, Link{
					Href: href,
					Text: text,
				})
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return links
}
