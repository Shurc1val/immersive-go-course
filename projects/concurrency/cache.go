package main

import (
	"fmt"
	"sync"
	"time"
)

type cacheStats struct {
	numWrites      int
	successfulHits int
	failedHits     int
	unreadDead     int
}

type keyHits[K comparable] struct {
	key  K
	hits int
}

type Cache[K comparable, V any] struct {
	mu          sync.RWMutex
	entryLimit  int
	data        map[K]V
	lru_tracker []keyHits[K]
	stats       cacheStats
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{
		entryLimit:  entryLimit,
		data:        map[K]V{},
		lru_tracker: []keyHits[K]{},
	}
}

// Removes an element from lru tracker, if said key exists; returns hit number from element.
func removeElement[K comparable](sl *[]keyHits[K], key K) int {
	for i, val := range *sl {
		if val.key == key {
			*sl = append((*sl)[:i], (*sl)[i+1:]...)
			return val.hits
		}
	}
	return 0
}

func (c *Cache[K, V]) resetLruTracker(key K) {
	hits := removeElement(&c.lru_tracker, key)
	c.lru_tracker = append(c.lru_tracker, keyHits[K]{key: key, hits: hits})
}

func (c *Cache[K, V]) enforceEntryLimit() bool {
	if len(c.lru_tracker) > c.entryLimit {
		killKey := c.lru_tracker[0].key
		if c.lru_tracker[0].hits == 0 {
			c.stats.unreadDead += 1
		}
		c.lru_tracker = c.lru_tracker[1:]
		delete(c.data, killKey)
		return true
	} else {
		return false
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether a value already existed in the cache for that key.
// If there was previously a value, it replaces that value with this one.
// Any Put counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.data[key]
	c.data[key] = value
	c.resetLruTracker(key)
	c.enforceEntryLimit()
	c.stats.numWrites += 1
	return ok
}

// Get returns the value assocated with the passed key, and a boolean to indicate whether a value was known or not. If not, nil is returned as the value.
// Any Get counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if ok {
		c.resetLruTracker(key)
		c.lru_tracker[len(c.lru_tracker)-1].hits += 1
		c.stats.successfulHits += 1
		return &val, true
	}
	c.stats.failedHits += 1
	return nil, false
}

func (c *Cache[K, V]) RetrieveStats() (hitRate float64, numNeverRead int, totalReadWrites int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hitRate = float64(c.stats.successfulHits) / float64(c.stats.successfulHits+c.stats.failedHits)

	count := 0
	for _, val := range c.lru_tracker {
		if val.hits == 0 {
			count += 1
		}
	}
	numNeverRead = count + c.stats.unreadDead

	totalReadWrites = c.stats.failedHits + c.stats.successfulHits + c.stats.numWrites

	return
}

func main() {
	c := NewCache[int, string](3)
	go c.Put(1, "Socks")
	go c.Put(2, "Hat")
	go c.Put(3, "Shoes")
	c.Get(2)
	go c.Put(4, "cloak")
	c.Get(5)
	time.Sleep(2 * time.Second)
	go fmt.Println(c.RetrieveStats())
	time.Sleep(3 * time.Second)
}
