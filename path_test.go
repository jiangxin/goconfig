package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	home, err = homeDir()
	assert.Nil(err)
	defer func(home string) {
		setHome(home)
	}(home)

	unsetHome()
	name, err = homeDir()
	assert.NotNil(err)
	assert.Equal("", name)

	name, err = expendHome("")
	assert.NotNil(err)
	assert.Equal("", name)

	setHome(tmpdir)

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

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = expendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = expendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

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

	home, err = homeDir()
	assert.Nil(err)
	defer func(home string) {
		setHome(home)
	}(home)

	unsetHome()
	name, err = absPath("~/")
	assert.NotNil(err)
	assert.Equal("", name)

	setHome(tmpdir)
	cwd, err := os.Getwd()
	assert.Nil(err)

	name, err = absPath("")
	assert.Nil(err, fmt.Sprintf("err should be nil, but got: %s", err))
	assert.Equal(cwd, name)

	name, err = absPath("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = absPath("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = absPath("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = absPath("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = absPath("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = absPath("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = absPath(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = absPath(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
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

	home, err = homeDir()
	assert.Nil(err)
	defer func(home string) {
		setHome(home)
	}(home)

	setHome(tmpdir)

	cwd := "/some/dir"
	if runtime.GOOS == "windows" {
		cwd = "c:\\some\\dir"
	}

	name, err = absJoin(cwd, "")
	assert.Nil(err)
	assert.Equal(cwd, name)

	name, err = absJoin(cwd, "a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = absJoin(cwd, "~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = absJoin(cwd, "~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = absJoin(cwd, "~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = absJoin(cwd, "~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = absJoin(cwd, "ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = absJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = absJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
}

func TestFindGitDir(t *testing.T) {
	var (
		err     error
		dir     string
		gitdir  string
		workdir string
		cfg     string
		home    string
		assert  = assert.New(t)
	)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = homeDir()
	assert.Nil(err)
	defer func(home string) {
		setHome(home)
	}(home)

	setHome(tmpdir)

	// find in: bare.git
	gitdir = filepath.Join(tmpdir, "bare.git")
	cmd := exec.Command("git", "init", "--bare", gitdir, "--")
	assert.Nil(cmd.Run())
	dir, err = findGitDir(gitdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in: bare.git/objects/pack
	dir, err = findGitDir(filepath.Join(gitdir, "objects", "pack"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(gitdir, "objects", "pack"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// create repo2 with gitdir file repo2/.git
	repo2 := filepath.Join(tmpdir, "repo2")
	err = os.MkdirAll(filepath.Join(repo2, "a", "b"), 0755)
	assert.Equal(nil, err)
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("gitdir: ../bare.git"),
		0644)
	assert.Nil(err)

	// find in: repo2/a/b/c
	dir, err = findGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// create bad gitdir file: repo2.git
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("../bare.git"),
		0644)
	assert.Nil(err)

	// fail to find in repo2/a/b/c (bad gitdir file)
	dir, err = findGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.NotNil(err)
	assert.Equal("", dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.NotNil(err)
	assert.Equal("", cfg)

	// create worktree
	workdir = filepath.Join(tmpdir, "workdir")
	cmd = exec.Command("git", "init", workdir, "--")
	assert.Nil(cmd.Run())

	gitdir = filepath.Join(workdir, ".git")
	err = os.MkdirAll(filepath.Join(workdir, "a", "b"), 0755)
	assert.Nil(err)

	// find in workdir
	dir, err = findGitDir(workdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(workdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = findGitDir(gitdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = findGitDir(filepath.Join(workdir, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(workdir, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// fail to find in tmpdir
	dir, err = findGitDir(tmpdir)
	assert.Equal("", dir)
	assert.Equal(ErrNotInGitDir, err)

	cfg, err = FindGitConfig(tmpdir)
	assert.Equal(ErrNotInGitDir, err)
	assert.Equal("", cfg)
}
