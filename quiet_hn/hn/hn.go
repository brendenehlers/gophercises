package hn

import (
	"log"
	"slices"
	"sync"
)

type Cache interface {
	Read(key string) any
	Set(key string, val any)
	Delete(key string)
}

func TopItems(count int) ([]Item, error) {
	c := client{}
	items := make([]Item, 0)
	ids, err := c.TopItems()
	if err != nil {
		return []Item{}, err
	}

	for offset := 0; len(items) < count; offset++ {
		start := offset * count
		end := (offset + 1) * count
		log.Printf("Processing ids %d to %d", start, end)
		batch := ids[start:end]

		itemsChan := make(chan Item, len(batch))
		var wg sync.WaitGroup

		for _, id := range batch {
			wg.Add(1)
			go worker(&wg, itemsChan, id)
		}

		wg.Wait()
		close(itemsChan)

		for item := range itemsChan {
			items = append(items, item)
			if len(items) >= count {
				break
			}
		}
	}

	slices.SortFunc(items, func(a Item, b Item) int {
		return a.ID - b.ID
	})

	return items, nil
}

func worker(wg *sync.WaitGroup, channel chan<- Item, id int) {
	defer wg.Done()

	c := client{}
	item, err := c.GetItem(id)
	if err != nil {
		panic(err)
	}

	if validate(item) {
		channel <- item
	}
}

func validate(item Item) bool {
	if item.Type == "story" && item.Url != "" {
		return true
	}
	return false
}
