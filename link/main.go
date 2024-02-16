package main

import (
	"bufio"
	"flag"
	"fmt"
	"gophercises/link/link"
	"os"
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

	links, err := link.Parse(reader)
	check(err)

	fmt.Println(links)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
