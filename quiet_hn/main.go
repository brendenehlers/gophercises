package main

import (
	"fmt"
	"gophercises/quiet_hn/hn"
)

func main() {
	var c hn.Client

	item, err := c.GetItem(39420256)
	if err != nil {
		panic(err)
	}

	fmt.Println(item)
}
