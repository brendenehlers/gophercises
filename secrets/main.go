package main

import (
	"fmt"
	"gophercises/secrets/encrypt"
	"gophercises/secrets/secret"
	"io"
	"os"
)

const (
	key = `6368616e676520746869732070617373`
)

func main() {
	v := secret.NewVault(key, "./secretsFile")
	err := v.Set("username", "myUsername")
	check(err)
	err = v.Set("password", "password")
	check(err)

	val, ok, err := v.Get("username")
	check(err)

	if ok {
		fmt.Println(val)
	} else {
		fmt.Println("value not found")
	}

	file, err := os.Open("./secretsFile")
	check(err)

	r, err := encrypt.DecryptReader([]byte(key), file)
	check(err)

	plaintext, err := io.ReadAll(r)
	check(err)

	fmt.Println("-----")
	fmt.Println(string(plaintext))
	fmt.Println("-----")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
