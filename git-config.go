package goconfig

import (
	"strconv"
	"strings"
)

// GitConfig maps section to key-value pairs
type GitConfig map[string]GitConfigKeys

// GitConfigKeys maps key to values
type GitConfigKeys map[string][]string

// NewGitConfig returns GitConfig with initialized maps
func NewGitConfig() GitConfig {
	c := make(GitConfig)
	return c
}

// Keys returns all config variable keys (in lower case)
func (v GitConfig) Keys() []string {
	allKeys := []string{}
	for s, keys := range v {
		for key := range keys {
			allKeys = append(allKeys, s+"."+key)
		}
	}
	return allKeys
}

// Add will add user input key-value pair
func (v GitConfig) Add(key, value string) {
	s, k := toSectionKey(key)
	v._add(s, k, value)
}

// _add key/value to config variables
func (v GitConfig) _add(section, key, value string) {
	// section, and key are always in lower case
	if _, ok := v[section]; !ok {
		v[section] = make(GitConfigKeys)
	}

	if _, ok := v[section][key]; !ok {
		v[section][key] = []string{}
	}
	v[section][key] = append(v[section][key], value)
}

// Get value from key
func (v GitConfig) Get(key string) string {
	values := v.GetAll(key)
	if values == nil || len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

// GetBool gets boolean from key with default value
func (v GitConfig) GetBool(key string, defaultValue bool) (bool, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	switch strings.ToLower(value) {
	case "yes", "true", "on":
		return true, nil
	case "no", "false", "off":
		return false, nil
	}
	return false, ErrNotBoolValue
}

// GetInt return integer value of key with default
func (v GitConfig) GetInt(key string, defaultValue int) (int, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

// GetInt64 return int64 value of key with default
func (v GitConfig) GetInt64(key string, defaultValue int64) (int64, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseInt(value, 10, 64)
}

// GetUint64 return uint64 value of key with default
func (v GitConfig) GetUint64(key string, defaultValue uint64) (uint64, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseUint(value, 10, 64)
}

// GetAll gets all values of a key
func (v GitConfig) GetAll(key string) []string {
	section, key := toSectionKey(key)

	keys := v[section]
	if keys != nil {
		return keys[key]
	}
	return nil
}

func dequoteKey(key string) string {
	if !strings.ContainsAny(key, "\"'") {
		return key
	}

	keys := []string{}
	for _, k := range strings.Split(key, ".") {
		keys = append(keys, strings.Trim(k, "\"'"))

	}
	return strings.Join(keys, ".")
}

// splitKey will split git config variable to section name and key
func toSectionKey(name string) (string, string) {
	name = strings.ToLower(dequoteKey(name))
	items := strings.Split(name, ".")

	if len(items) < 2 {
		return "", ""
	}
	key := items[len(items)-1]
	section := strings.Join(items[0:len(items)-1], ".")
	return section, key
}
