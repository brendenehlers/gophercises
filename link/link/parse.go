package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	var links []Link
	for _, n := range nodes {
		links = append(links, buildLink(n))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var href string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	text := text(n)

	return Link{
		Href: href,
		Text: text,
	}
}

func text(n *html.Node) string {
	var t string
	if n.Type == html.TextNode {
		t += n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		t += text(c)
	}

	return strings.Join(strings.Fields(t), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n} // [n] equivalent
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}
