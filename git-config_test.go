package goconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidSectionName(t *testing.T) {
	assert := assert.New(t)

	data := `# The following section name should have quote, like: [a "b"]
[a b]
	c = d`
	_, lineno, err := Parse([]byte(data), "filename")
	assert.Equal(ErrMissingStartQuote, err)
	assert.Equal(uint(2), lineno)
}

func TestInvalidKeyWithSpace(t *testing.T) {
	assert := assert.New(t)

	data := `# keys should not have spaces
[a]
	b c = d`
	_, lineno, err := Parse([]byte(data), "filename")
	assert.Equal(ErrInvalidKeyChar, err)
	assert.Equal(uint(3), lineno)
}

func TestParseSectionWithSpaces1(t *testing.T) {
	assert := assert.New(t)

	data := `[ab "cd"]
	value1 = x
	value2 = x y
	value3  = a \"quote
[remote "hello world"]
	url = test`
	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)
	assert.Equal("x", cfg.Get("ab.cd.value1"))
	assert.Equal("x y", cfg.Get("ab.cd.value2"))
	assert.Equal("a \"quote", cfg.Get("ab.cd.value3"))
}

func TestParseSectionWithSpaces2(t *testing.T) {
	assert := assert.New(t)

	data := `[remote "hello world"]
	url = test`
	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)
	assert.Equal("test", cfg.Get("remote.hello world.url"))
	assert.Equal("test", cfg.Get(`remote."hello world".url`))
	assert.Equal("test", cfg.Get(`"remote.hello world".url`))
	assert.Equal("test", cfg.Get(`"remote.hello world.url"`))
}

func TestGetAll(t *testing.T) {
	assert := assert.New(t)

	data := `[remote "origin"]
	url = https://example.com/my/repo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
	fetch = +refs/tags/*:refs/tags/*`
	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)
	assert.Equal("+refs/tags/*:refs/tags/*", cfg.Get("remote.origin.fetch"))
	assert.Equal([]string{
		"+refs/heads/*:refs/remotes/origin/*",
		"+refs/tags/*:refs/tags/*",
	}, cfg.GetAll("remote.origin.fetch"))

}

func TestGetBool(t *testing.T) {
	assert := assert.New(t)

	data := `[a]
	t1 = true
	t2 = yes
	t3 = on
	f1 = false
	f2 = no
	f3 = off
	x1 = 1
	x2 = nothing`

	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)

	v, err := cfg.GetBool("a.t1", false)
	assert.Nil(err)
	assert.True(v)

	v, err = cfg.GetBool("a.t2", false)
	assert.Nil(err)
	assert.True(v)

	v, err = cfg.GetBool("a.t3", false)
	assert.Nil(err)
	assert.True(v)

	v, err = cfg.GetBool("a.t4", false)
	assert.Nil(err)
	assert.False(v)

	v, err = cfg.GetBool("a.f1", true)
	assert.Nil(err)
	assert.False(v)

	v, err = cfg.GetBool("a.f2", true)
	assert.Nil(err)
	assert.False(v)

	v, err = cfg.GetBool("a.f3", true)
	assert.Nil(err)
	assert.False(v)

	v, err = cfg.GetBool("a.f4", true)
	assert.Nil(err)
	assert.True(v)

	v, err = cfg.GetBool("a.x1", true)
	assert.Equal(ErrNotBoolValue, err)

	v, err = cfg.GetBool("a.x2", true)
	assert.Equal(ErrNotBoolValue, err)
}

func TestGetInt(t *testing.T) {
	assert := assert.New(t)

	data := `[a]
	i1 = 1
	i2 = 100
	i3 = abc`

	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)

	v1, err := cfg.GetInt("a.i1", 0)
	assert.Nil(err)
	assert.Equal(1, v1)

	v2, err := cfg.GetInt64("a.i2", 0)
	assert.Nil(err)
	assert.Equal(int64(100), v2)

	v3, err := cfg.GetUint64("a.i2", 0)
	assert.Nil(err)
	assert.Equal(uint64(100), v3)

	_, err = cfg.GetInt("a.i3", 0)
	assert.NotNil(err)

	v4, err := cfg.GetInt("a.i4", 6700)
	assert.Nil(err)
	assert.Equal(6700, v4)
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	data := `[a]
	b = value-b
	c = value-c`

	cfg, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)

	assert.Equal("value-b", cfg.Get("a.b"))
	assert.Equal("value-c", cfg.Get("a.c"))

	data = `[a]
	c = other-c
	d = other-d`

	cfg2, _, err := Parse([]byte(data), "filename")
	assert.Nil(err)
	assert.Equal("other-c", cfg2.Get("a.c"))
	assert.Equal("other-d", cfg2.Get("a.d"))

	cfg.Merge(cfg2)
	assert.Equal("value-b", cfg.Get("a.b"))
	assert.Equal("other-c", cfg.Get("a.c"))
	assert.Equal("other-d", cfg.Get("a.d"))
	assert.Equal([]string{
		"value-c",
		"other-c",
	}, cfg.GetAll("a.c"))
}
