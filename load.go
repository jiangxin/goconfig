package goconfig

import (
	"io/ioutil"
	"os"
)

// Load only loads one file or config of current repository
func Load(name string) (GitConfig, error) {
	var (
		err    error
		fi     os.FileInfo
		search bool
	)

	if name == "" {
		search = true
		name, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	} else {
		fi, err = os.Stat(name)
		if err != nil {
			return nil, ErrNotExist
		}
		if fi.IsDir() {
			search = true
		}
	}

	if search {
		name, err = FindGitConfig(name)
		if err != nil {
			if err == ErrNotInGitDir {
				return nil, nil
			}
			return nil, err
		}
	}

	c, ok := configCaches.get(name)
	if ok {
		return c.config, nil
	}

	// cache will be updated using this time
	fi, _ = os.Stat(name)

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	cfg, _, err := Parse(buf, name)
	if err != nil {
		return cfg, err
	}

	// update cache
	configCaches.set(name, cfg, fi.ModTime())
	return cfg, nil
}

// LoadAll will load additional global and system config files
func LoadAll(name string) (GitConfig, error) {
	cfg := NewGitConfig()

	sysConfig, err := SystemConfig()
	if err != nil {
		return nil, err
	}

	globalConfig, err := GlobalConfig()
	if err != nil {
		return nil, err
	}

	repoConfig, err := Load(name)
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

// SystemConfig returns system git config, reload if necessary
func SystemConfig() (GitConfig, error) {
	file := SystemConfigFile()
	if file == "" {
		return nil, nil
	}

	if _, err := os.Stat(file); err != nil {
		return nil, nil
	}
	return Load(file)
}

// GlobalConfig returns global user config, reload if necessary
func GlobalConfig() (GitConfig, error) {
	file, err := GlobalConfigFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(file); err != nil {
		return nil, nil
	}
	return Load(file)
}
