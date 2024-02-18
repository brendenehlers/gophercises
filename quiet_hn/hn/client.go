package hn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	Text string `json:"text"`
	Url  string `json:"url"`
}

type client struct {
	apiBase string
}

func (c *client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}
}

func (c *client) get(context string) (*http.Response, error) {
	log.Printf("Requesting: %s\n", context)
	resp, err := http.Get(c.apiBase + context)
	return resp, err
}

func (c *client) TopItems() ([]int, error) {
	c.defaultify()
	resp, err := c.get("/topstories.json")
	if err != nil {
		return []int{}, err
	}
	defer resp.Body.Close()

	var ids []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return []int{}, err
	}

	return ids, nil
}

func (c *client) GetItem(id int) (Item, error) {
	c.defaultify()

	resp, err := c.get(fmt.Sprintf("/item/%d.json", id))
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	var item Item
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}
