package main

import (
	"fmt"
	"gophercises/secrets/encrypt"
)

const (
	secret = `6368616e676520746869732070617373`
	// ciphertext = `cf0495cc6f75dafc23948538e79904a9`
	plaintext = `My super secret text`
)

func main() {

	ciphertext, err := encrypt.Encrypt([]byte(secret), []byte(plaintext))
	check(err)

	text, err := encrypt.Decrypt([]byte(secret), ciphertext)
	check(err)

	fmt.Printf("%s\n", text)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
