package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
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

func loadOne(name string, searching bool) (GitConfig, error) {
	var (
		err error
		fi  os.FileInfo
	)

	if name == "" {
		if !searching {
			return nil, fmt.Errorf("filepath is not provided")
		}

		name, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	// cannot find, lookup upper dir to find gitdir
	fi, err = os.Stat(name)
	if err != nil || fi.IsDir() {
		if !searching {
			return nil, nil
		}

		name, err = FindGitConfig(name)
		if err != nil {
			if err == ErrNotInGitDir {
				return nil, nil
			} else {
				return nil, err
			}
		}

		fi, err = os.Stat(name)

		if err != nil {
			return nil, nil
		}
	}

	// return cache if available
	c, ok := configCaches.Get(name)
	if ok {
		if fi.ModTime() == c.time {
			return c.config, nil
		}
	}

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	cfg, _, err := Parse(buf, name)
	if err != nil {
		return cfg, err
	}

	// update cache
	configCaches.Set(name, cacheItem{
		config: cfg,
		time:   fi.ModTime(),
	})
	return cfg, nil
}

func loadAll(name string, searching bool) (GitConfig, error) {
	cfg := NewGitConfig()

	sysConfig, err := SystemConfig()
	if err != nil {
		return nil, err
	}

	globalConfig, err := GlobalConfig()
	if err != nil {
		return nil, err
	}

	repoConfig, err := loadOne(name, searching)
	if err != nil {
		return nil, err
	}

	if sysConfig != nil {
		cfg.Merge(sysConfig, ScopeSystem)
	}

	if globalConfig != nil {
		cfg.Merge(globalConfig, ScopeGlobal)
	}

	if repoConfig != nil {
		cfg.Merge(repoConfig, ScopeSelf)
	}

	return cfg, nil
}

// LoadFile loads git config file
func LoadFile(name string, inherit bool) (GitConfig, error) {
	if !inherit {
		return loadOne(name, false)
	}

	return loadAll(name, false)
}

// LoadDir loads config from gitdir
func LoadDir(name string, inherit bool) (GitConfig, error) {
	if !inherit {
		return loadOne(name, true)
	}

	return loadAll(name, true)
}

// SystemConfig returns system git config, reload if necessary
func SystemConfig() (GitConfig, error) {
	file := SystemConfigFile()
	if file == "" {
		return nil, nil
	}

	return loadOne(file, false)
}

// GlobalConfig returns global user config, reload if necessary
func GlobalConfig() (GitConfig, error) {
	file, err := GlobalConfigFile()
	if err != nil {
		return nil, err
	}

	return loadOne(file, false)
}
