package main

import (
	"fmt"
	"gophercises/quiet_hn/client"
)

func main() {
	var c client.Client

	item, err := c.GetItem(39420256)
	if err != nil {
		panic(err)
	}

	fmt.Println(item)
}
