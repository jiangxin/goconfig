package goconfig

import (
	"io/ioutil"
	"os"
	"os/exec"
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

func TestFindGitDir(t *testing.T) {
	var (
		err     error
		dir     string
		gitdir  string
		workdir string
		cfg     string
	)

	home := os.Getenv("HOME")

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	os.Setenv("HOME", tmpdir)

	// find in: bare.git
	gitdir = filepath.Join(tmpdir, "bare.git")
	cmd := exec.Command("git", "init", "--bare", gitdir, "--")
	assert.Equal(t, nil, cmd.Run())
	dir, err = FindGitDir(gitdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// find in: bare.git/objects/pack
	dir, err = FindGitDir(filepath.Join(gitdir, "objects", "pack"))
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(gitdir, "objects", "pack"))
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// create repo2 with gitdir file repo2/.git
	repo2 := filepath.Join(tmpdir, "repo2")
	err = os.MkdirAll(filepath.Join(repo2, "a", "b"), 0755)
	assert.Equal(t, nil, err)
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("gitdir: ../bare.git"),
		0644)
	assert.Equal(t, nil, err)

	// find in: repo2/a/b/c
	dir, err = FindGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// create bad gitdir file: repo2.git
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("../bare.git"),
		0644)
	assert.Equal(t, nil, err)

	// fail to find in repo2/a/b/c (bad gitdir file)
	dir, err = FindGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", cfg)

	// create worktree
	workdir = filepath.Join(tmpdir, "workdir")
	cmd = exec.Command("git", "init", workdir, "--")
	assert.Equal(t, nil, cmd.Run())

	gitdir = filepath.Join(workdir, ".git")
	err = os.MkdirAll(filepath.Join(workdir, "a", "b"), 0755)
	assert.Equal(t, nil, err)

	// find in workdir
	dir, err = FindGitDir(workdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(workdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = FindGitDir(gitdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = FindGitDir(filepath.Join(workdir, "a", "b", "c"))
	assert.Equal(t, nil, err)
	assert.Equal(t, gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(workdir, "a", "b", "c"))
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(gitdir, "config"), cfg)

	// fail to find in tmpdir
	dir, err = FindGitDir(tmpdir)
	assert.Equal(t, "", dir)
	assert.Equal(t, nil, err)

	cfg, err = FindGitConfig(tmpdir)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", cfg)

	os.Setenv("HOME", home)
}
