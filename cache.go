package goconfig

import (
	"sync"
	"time"
)

type cache struct {
	caches map[string]cacheItem
	mu     sync.RWMutex
}

type cacheItem struct {
	config GitConfig
	time   time.Time
}

var (
	configCaches = cache{caches: make(map[string]cacheItem)}
)

func (v *cache) Set(key string, item cacheItem) {
	v.mu.Lock()
	v.caches[key] = item
	v.mu.Unlock()
}

func (v *cache) Get(key string) (cacheItem, bool) {
	v.mu.RLock()
	item, ok := v.caches[key]
	v.mu.RUnlock()
	return item, ok
}
