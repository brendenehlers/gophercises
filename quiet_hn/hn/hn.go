package hn

import (
	"log"
	"slices"
	"sync"
)

type Cache interface {
	Read(key int) (Item, bool)
	Insert(key int, val Item) error
	Remove(key int) error
}

type HN struct {
	Cache Cache
}

func (hn *HN) TopItems(count int) ([]Item, error) {
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
			if val, ok := hn.Cache.Read(id); ok {
				itemsChan <- val
			} else {
				wg.Add(1)
				go hn.worker(&wg, itemsChan, id)
			}
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

func (hn *HN) worker(wg *sync.WaitGroup, channel chan<- Item, id int) {
	defer wg.Done()

	c := client{}
	item, err := c.GetItem(id)
	if err != nil {
		panic(err)
	}

	if hn.Cache != nil {
		hn.Cache.Insert(id, item)
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
