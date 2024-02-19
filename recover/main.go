package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	rw.writes = append(rw.writes, bytes)
	return len(bytes), nil
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}

	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijacker not supported")
	}

	return hj.Hijack()
}

func (rw *responseWriter) Flush() {
	f, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	f.Flush()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)

	log.Println("Hosted server on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", recoverMW(mux, true)))
}

func recoverMW(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered error:", r)

				if dev {
					log.Println("Panic stack trace:")
					stack := string(debug.Stack())
					log.Println(stack)
				}

				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
		}()
		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
