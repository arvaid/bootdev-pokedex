package internal

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex   *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	var cache Cache
	cache.entries = make(map[string]cacheEntry)
	cache.mutex = &sync.Mutex{}
	go cache.reapLoop(interval)
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if entry, ok := cache.entries[key]; ok {
		return entry.val, true
	}
	return nil, false
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			cache.mutex.Lock()
			for key, entry := range cache.entries {
				if entry.createdAt.Add(interval).Compare(t) >= 0 {
					delete(cache.entries, key)
				}
			}
			cache.mutex.Unlock()
		}
	}
}
