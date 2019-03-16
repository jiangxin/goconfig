package goconfig

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemConfig(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	cfgFile := filepath.Join(tmpdir, "test.config")
	os.Setenv(gitSystemConfigEnv, cfgFile)
	defer os.Unsetenv(gitSystemConfigEnv)

	err = ioutil.WriteFile(cfgFile,
		[]byte(`[test]
	foo = system foo`),
		0644)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(cfgFile)
	if err != nil {
		panic(err)
	}

	cacheTime := time.Now().Add(time.Duration(-999) * time.Second)

	// reset timestamp of system config file
	os.Chtimes(cfgFile, cacheTime, cacheTime)

	cfg, err := SystemConfig()
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal("system foo", cfg.Get("test.foo"))

	// sysconfig changed
	err = ioutil.WriteFile(cfgFile,
		[]byte(`[test]
	foo = system foobar`),
		0644)
	if err != nil {
		panic(err)
	}

	// reset timestamp of system config file
	os.Chtimes(cfgFile, cacheTime, cacheTime)

	// using cache
	cfg, err = SystemConfig()
	assert.NotNil(cfg)
	assert.Equal("system foo", cfg.Get("test.foo"))

	// change sysconfig again with new timestamp
	err = ioutil.WriteFile(cfgFile,
		[]byte(`[test]
	foo = system foobaz`),
		0644)
	if err != nil {
		panic(err)
	}

	// timestamp changed, auto reload
	cfg, err = SystemConfig()
	assert.NotNil(cfg)
	assert.Equal("system foobaz", cfg.Get("test.foo"))
}

func TestLoadGlobalConfig(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)
	defer os.Setenv("HOME", home)

	cfg, err := GlobalConfig()
	assert.Nil(err)
	assert.Nil(cfg)

	cfgFile := filepath.Join(tmpdir, ".gitconfig")
	err = ioutil.WriteFile(cfgFile,
		[]byte(`[user]
	name = user1
	email = user1@email.addr`),
		0644)
	if err != nil {
		panic(err)
	}

	cacheTime := time.Now().Add(time.Duration(-999) * time.Second)
	// reset timestamp of system config file
	os.Chtimes(cfgFile, cacheTime, cacheTime)

	cfg, err = GlobalConfig()
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal("user1", cfg.Get("user.name"))
	assert.Equal("user1@email.addr", cfg.Get("user.email"))

	err = ioutil.WriteFile(cfgFile,
		[]byte(`[user]
	name = user2
	email = user2@email.addr`),
		0644)
	if err != nil {
		panic(err)
	}

	// reset timestamp of system config file
	os.Chtimes(cfgFile, cacheTime, cacheTime)

	// using cache
	cfg, err = GlobalConfig()
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal("user1", cfg.Get("user.name"))
	assert.Equal("user1@email.addr", cfg.Get("user.email"))

	err = ioutil.WriteFile(cfgFile,
		[]byte(`[user]
	name = user3
	email = user3@email.addr`),
		0644)
	if err != nil {
		panic(err)
	}

	// timestamp of changed, refresh cache
	cfg, err = GlobalConfig()
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal("user3", cfg.Get("user.name"))
	assert.Equal("user3@email.addr", cfg.Get("user.email"))
}

func TestRepoConfig(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)
	defer os.Setenv("HOME", home)

	// create bare.git
	gitdir := filepath.Join(tmpdir, "bare.git")
	cmd := exec.Command("git", "init", "--bare", gitdir, "--")
	assert.Nil(cmd.Run())

	// create shared config
	sharedCfg := filepath.Join(tmpdir, "shared.config")
	cmd = exec.Command("git", "config", "-f", sharedCfg, "test.foo", "has foo")
	assert.Nil(cmd.Run())

	// set config in bare.git
	cmd = exec.Command("git", "-C", gitdir, "config", "include.path", sharedCfg)
	assert.Nil(cmd.Run())
	cmd = exec.Command("git", "-C", gitdir, "config", "test.bar", "has bar")
	assert.Nil(cmd.Run())

	// load config of bare.git
	cfg, err := LoadDir(gitdir, false)
	assert.Nil(err)
	assert.Equal("has bar", cfg.Get("test.bar"))
	assert.Equal([]string{sharedCfg}, cfg.GetAll("include.path"))
	assert.Equal("has foo", cfg.Get("test.foo"))
}

func TestCircularInclude(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	// create circular include
	sharedCfg1 := filepath.Join(tmpdir, "shared.config.1")
	sharedCfg2 := filepath.Join(tmpdir, "shared.config.2")
	cmd := exec.Command("git", "config", "-f", sharedCfg1, "include.path", sharedCfg2)
	assert.Nil(cmd.Run())
	cmd = exec.Command("git", "config", "-f", sharedCfg2, "include.path", sharedCfg1)
	assert.Nil(cmd.Run())

	_, err = LoadFile(sharedCfg1, false)
	assert.NotNil(err)
	assert.True(strings.HasPrefix(err.Error(), "exceeded maximum include depth"))

	// include circular include in test repo
	workdir := filepath.Join(tmpdir, "workdir")
	cmd = exec.Command("git", "init", workdir, "--")
	assert.Nil(cmd.Run())
	cmd = exec.Command("git", "-C", workdir, "config", "include.path", sharedCfg1)
	assert.Nil(cmd.Run())
	_, err = LoadDir(workdir, false)
	assert.True(strings.HasPrefix(err.Error(), "exceeded maximum include depth"))
}

func TestAllConfig(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	// Create system config
	sysCfgFile := filepath.Join(tmpdir, "system-config")
	os.Setenv(gitSystemConfigEnv, sysCfgFile)
	defer os.Unsetenv(gitSystemConfigEnv)

	assert.Nil(exec.Command("git", "config", "-f", sysCfgFile, "test.key1", "sys 1").Run())
	assert.Nil(exec.Command("git", "config", "-f", sysCfgFile, "test.key2", "sys 2").Run())
	assert.Nil(exec.Command("git", "config", "-f", sysCfgFile, "test.key3", "sys 3").Run())
	sysConfig, err := SystemConfig()
	assert.Nil(err)
	assert.NotNil(sysConfig)

	// Create user config
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)
	defer os.Setenv("HOME", home)

	userCfgFile, err := GlobalConfigFile()
	assert.Nil(err)
	assert.Nil(exec.Command("git", "config", "-f", userCfgFile, "test.key2", "user 2").Run())
	assert.Nil(exec.Command("git", "config", "-f", userCfgFile, "test.key3", "user 3").Run())
	assert.Nil(exec.Command("git", "config", "-f", userCfgFile, "test.key4", "user 4").Run())
	userConfig, err := GlobalConfig()
	assert.Nil(err)
	assert.NotNil(userConfig)

	// Create repo config
	workdir := filepath.Join(tmpdir, "workdir")
	assert.Nil(exec.Command("git", "init", workdir, "--").Run())
	assert.Nil(exec.Command("git", "-C", workdir, "config", "test.key3", "repo 3").Run())
	assert.Nil(exec.Command("git", "-C", workdir, "config", "test.key4", "repo 4").Run())
	assert.Nil(exec.Command("git", "-C", workdir, "config", "test.key5", "repo 5").Run())

	repoConfig, err := LoadDir(workdir, false)
	assert.Nil(err)
	assert.NotNil(repoConfig)

	// Get all config
	allConfig, err := LoadDir(workdir, true)
	assert.Nil(err)
	assert.NotNil(allConfig)

	// Check system config
	assert.Equal([]string{"sys 1"}, sysConfig.GetAll("test.key1"))
	assert.Equal([]string{"sys 2"}, sysConfig.GetAll("test.key2"))
	assert.Equal([]string{"sys 3"}, sysConfig.GetAll("test.key3"))
	assert.Equal(true, nil == sysConfig.GetAll("test.key4"))
	assert.Equal(true, nil == sysConfig.GetAll("test.key5"))

	// Check global config
	assert.Equal(true, nil == userConfig.GetAll("test.key1"))
	assert.Equal([]string{"user 2"}, userConfig.GetAll("test.key2"))
	assert.Equal([]string{"user 3"}, userConfig.GetAll("test.key3"))
	assert.Equal([]string{"user 4"}, userConfig.GetAll("test.key4"))
	assert.Equal(true, nil == userConfig.GetAll("test.key5"))

	// Check repo config
	assert.Equal(true, nil == repoConfig.GetAll("test.key1"))
	assert.Equal(true, nil == repoConfig.GetAll("test.key2"))
	assert.Equal([]string{"repo 3"}, repoConfig.GetAll("test.key3"))
	assert.Equal([]string{"repo 4"}, repoConfig.GetAll("test.key4"))
	assert.Equal([]string{"repo 5"}, repoConfig.GetAll("test.key5"))

	// Check merged config
	assert.Equal([]string{"sys 1"}, allConfig.GetAll("test.key1"))
	assert.Equal([]string{"sys 2", "user 2"}, allConfig.GetAll("test.key2"))
	assert.Equal([]string{"sys 3", "user 3", "repo 3"}, allConfig.GetAll("test.key3"))
	assert.Equal([]string{"user 4", "repo 4"}, allConfig.GetAll("test.key4"))
	assert.Equal([]string{"repo 5"}, allConfig.GetAll("test.key5"))

	assert.Equal("sys 1", allConfig.Get("test.key1"))
	assert.Equal("user 2", allConfig.Get("test.key2"))
	assert.Equal("repo 3", allConfig.Get("test.key3"))
	assert.Equal("repo 4", allConfig.Get("test.key4"))
	assert.Equal("repo 5", allConfig.Get("test.key5"))
}

func TestSaveConfig(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)
	defer os.Setenv("HOME", home)

	// Create system config
	cfgFile := filepath.Join(tmpdir, "gitconfig")

	assert.Nil(exec.Command("git", "config", "-f", cfgFile, "ab.cd.ef", "value-1").Run())
	assert.Nil(exec.Command("git", "config", "-f", cfgFile, "ab.cd e.fg", "value 2").Run())
	assert.Nil(exec.Command("git", "config", "-f", cfgFile, "--add", "ab.cd e.fg", "value 3").Run())
	assert.Nil(exec.Command("git", "config", "-f", cfgFile, "--add", "ab.cd e.fg", "value 4").Run())
	assert.Nil(exec.Command("git", "config", "-f", cfgFile, "ab.cd", "value has space ").Run())

	// Load cfgFile
	cfg, err := LoadFile(cfgFile, false)
	assert.Nil(err)
	assert.Equal("value has space ", cfg.Get("ab.cd"))

	// Save file
	newCfgFile := cfgFile + ".new"
	err = cfg.Save(newCfgFile)
	assert.Nil(err)
	_, err = os.Stat(newCfgFile)
	assert.Nil(err)

	// Load new file
	newCfg, err := LoadFile(newCfgFile, false)
	assert.Nil(err)
	assert.Equal(cfg, newCfg)
}
