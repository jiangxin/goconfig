package goconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpendHome(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home = os.Getenv("HOME")

	os.Unsetenv("HOME")
	name, err = homeDir()
	assert.NotNil(err)
	assert.Equal("", name)

	name, err = expendHome("")
	assert.NotNil(err)
	assert.Equal("", name)

	os.Setenv("HOME", tmpdir)

	name, err = homeDir()
	assert.Equal(tmpdir, name)

	name, err = expendHome("")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = expendHome("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = expendHome("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "~a"), name)

	name, err = expendHome("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = expendHome("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = expendHome("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = expendHome("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "ab"), name)

	name, err = expendHome("/")
	assert.Nil(err)
	assert.Equal("/", name)

	name, err = expendHome("/a")
	assert.Nil(err)
	assert.Equal("/a", name)

	os.Setenv("HOME", home)
}

func TestAbs(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home = os.Getenv("HOME")

	os.Unsetenv("HOME")
	name, err = Abs("~/")
	assert.NotNil(err)
	assert.Equal("", name)

	os.Setenv("HOME", tmpdir)
	cwd, _ := os.Getwd()

	name, err = Abs("")
	assert.Nil(err)
	assert.Equal(cwd, name)

	name, err = Abs("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = Abs("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = Abs("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = Abs("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	name, err = Abs("/")
	assert.Nil(err)
	assert.Equal("/", name)

	name, err = Abs("/a")
	assert.Nil(err)
	assert.Equal("/a", name)

	os.Setenv("HOME", home)
}

func TestAbsJoin(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home = os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)

	cwd := "/some/dir"

	name, err = AbsJoin(cwd, "")
	assert.Nil(err)
	assert.Equal(cwd, name)

	name, err = AbsJoin(cwd, "a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = AbsJoin(cwd, "~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = AbsJoin(cwd, "~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = AbsJoin(cwd, "ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	name, err = AbsJoin(cwd, "/")
	assert.Nil(err)
	assert.Equal("/", name)

	name, err = AbsJoin(cwd, "/a")
	assert.Nil(err)
	assert.Equal("/a", name)

	os.Setenv("HOME", home)
}
