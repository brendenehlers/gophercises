package main

import (
	"flag"
	"fmt"
)

func main() {
	filepathFlag := flag.String("file", "samples/ex1.html", "file to read")
	flag.Parse()

	fmt.Println("file:", *filepathFlag)

}
