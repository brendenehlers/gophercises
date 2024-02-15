package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	yamlFile := flag.String("yaml", "", "provides a yaml file to read redirect paths from")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	// pathsToUrls := map[string]string{
	// 	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	// 	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	// 	"/google":         "https://google.com",
	// }
	mapHandler := urlshort.MapHandler(make(map[string]string), mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yamlString := `
	// - route:
	//     path: /urlshort
	//     url: https://github.com/gophercises/urlshort
	// - route:
	//     path: /urlshort-final
	//     url: https://github.com/gophercises/urlshort/tree/solution
	// `

	data, err := os.ReadFile(*yamlFile)
	if err != nil {
		panic(err)
	}
	yaml := data

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
