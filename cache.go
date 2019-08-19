package goconfig

import (
	"os"
	"sync"
	"time"
)

type cache struct {
	caches map[string]*cacheItem
	mu     sync.RWMutex
}

type cacheItem struct {
	config   GitConfig
	filename string
	time     time.Time
	size     int64
}

var (
	configCaches = cache{caches: make(map[string]*cacheItem)}
)

func (v *cacheItem) uptodate() bool {
	fi, err := os.Stat(v.filename)
	if err == nil && fi.ModTime() == v.time && fi.Size() == v.size {
		return true
	}
	return false
}

// set will cache cache entry
func (v *cache) set(key string, cfg GitConfig, size int64, modTime time.Time) {
	v.mu.Lock()
	v.caches[key] = &cacheItem{
		config:   cfg,
		filename: key,
		time:     modTime,
		size:     size,
	}
	v.mu.Unlock()
}

// get returns cache entry if uptodate
func (v *cache) get(key string) (*cacheItem, bool) {
	v.mu.RLock()
	item, ok := v.caches[key]
	v.mu.RUnlock()

	if !ok {
		return nil, false
	}
	if !item.uptodate() {
		v.expire(key)
		return nil, false
	}
	return item, ok
}

func (v *cache) expire(key string) {
	v.mu.Lock()
	delete(v.caches, key)
	v.mu.Unlock()
}
