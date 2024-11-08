package pokeCache

import (
	"fmt"
	"sync"
	"time"
)

type PokeCache struct {
	cacheEntry map[string]cacheEntry
	mutex      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(d time.Duration) *PokeCache {
	ce := make(map[string]cacheEntry)
	c := PokeCache{cacheEntry: ce}
	go c.ReapLoop(d)
	return &c
}

func (c *PokeCache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	newEntry := cacheEntry{val: val, createdAt: time.Now()}
	c.cacheEntry[key] = newEntry
}

func (c *PokeCache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exists := c.cacheEntry[key]
	if !exists {
		return []byte{}, false
	}
	fmt.Println("value found!")
	return entry.val, true
}

func (c *PokeCache) ReapLoop(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		timeup := <-ticker.C
		//fmt.Println("reaped at time: " + timeup.String())
		c.mutex.Lock()
		for key, ele := range c.cacheEntry {
			if ele.createdAt.Before(timeup) {
				delete(c.cacheEntry, key)
			}
		}
		c.mutex.Unlock()
	}
}
