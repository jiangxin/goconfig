package goconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	gitSystemConfigEnv = "TEST_GIT_SYSTEM_CONFIG"
)

// homeDir returns home directory
func homeDir() (string, error) {
	var (
		home string
	)

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	}
	if home == "" {
		home = os.Getenv("HOME")
	}

	if home == "" {
		return "", fmt.Errorf("cannot find HOME")
	}

	return home, nil
}

func xdgConfigHome(file string) (string, error) {
	var (
		home string
		err  error
	)

	home = os.Getenv("XDG_CONFIG_HOME")
	if home != "" {
		return filepath.Join(home, "git", file), nil
	}

	home, err = homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "git", file), nil
}

// expendHome expends path prefix "~/" to home dir
func expendHome(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	home, err := homeDir()
	if err != nil {
		return "", err
	}

	if len(name) == 0 || name == "~" {
		return home, nil
	} else if len(name) > 1 && name[0] == '~' && (name[1] == '/' || name[1] == '\\') {
		return filepath.Join(home, name[2:]), nil
	}

	return filepath.Join(home, name), nil
}

// absPath returns absolute path and will expend homedir if path has "~/' prefix
func absPath(name string) (string, error) {
	if name == "" {
		return os.Getwd()
	}

	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return expendHome(name)
	}

	return filepath.Abs(name)
}

// absJoin returns absolute path, and use <dir> as parent dir for relative path
func absJoin(dir, name string) (string, error) {
	if name == "" {
		return filepath.Abs(dir)
	}

	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return expendHome(name)
	}

	return absPath(filepath.Join(dir, name))
}

// isGitDir test whether dir is a valid git dir
func isGitDir(dir string) bool {
	var (
		err error
		fi  os.FileInfo
	)

	objectDir := filepath.Join(dir, "objects", "pack")
	if fi, err = os.Stat(objectDir); err != nil || !fi.IsDir() {
		return false
	}

	refsDir := filepath.Join(dir, "refs")
	if fi, err = os.Stat(refsDir); err != nil || !fi.IsDir() {
		return false
	}

	cfgFile := filepath.Join(dir, "config")
	if fi, err = os.Stat(cfgFile); err != nil || fi.IsDir() {
		return false
	}

	return true
}

// findGitDir searches git dir
func findGitDir(dir string) (string, error) {
	var err error

	dir, err = absPath(dir)
	if err != nil {
		return "", err
	}

	for {
		// Check if is in a bare repo
		if isGitDir(dir) {
			return dir, nil
		}

		// Check .git
		gitdir := filepath.Join(dir, ".git")
		fi, err := os.Stat(gitdir)
		if err != nil {
			// Test parent dir
			oldDir := dir
			dir = filepath.Dir(dir)
			if oldDir == dir {
				break
			}
			continue
		} else if fi.IsDir() {
			if isGitDir(gitdir) {
				return gitdir, nil
			}
			return "", fmt.Errorf("corrupt git dir: %s", gitdir)
		} else {
			f, err := os.Open(gitdir)
			if err != nil {
				return "", fmt.Errorf("cannot open gitdir file '%s'", gitdir)
			}
			defer f.Close()
			reader := bufio.NewReader(f)
			line, err := reader.ReadString('\n')
			if strings.HasPrefix(line, "gitdir:") {
				realgit := strings.TrimSpace(strings.TrimPrefix(line, "gitdir:"))
				if !filepath.IsAbs(realgit) {
					realgit, err = absJoin(filepath.Dir(gitdir), realgit)
					if err != nil {
						return "", err
					}
				}
				if isGitDir(realgit) {
					return realgit, nil
				}
				return "", fmt.Errorf("gitdir '%s' points to corrupt git repo: %s", gitdir, realgit)
			}
			return "", fmt.Errorf("bad gitdir file '%s'", gitdir)
		}
	}
	return "", ErrNotInGitDir
}

// FindGitConfig returns local git config file
func FindGitConfig(dir string) (string, error) {
	dir, err := findGitDir(dir)
	if err == nil {
		return filepath.Join(dir, "config"), nil
	}
	return "", err
}

// SystemConfigFile returns system git config file
func SystemConfigFile() string {
	file := os.Getenv(gitSystemConfigEnv)
	if file == "" {
		file = "/etc/gitconfig"
	}
	return file
}

// GlobalConfigFile returns global git config file
func GlobalConfigFile() (string, error) {
	var (
		file string
		err  error
	)

	file, err = xdgConfigHome("config")
	if err != nil {
		return "", err
	}

	// xdg config not exist, use ~/.gitconfig
	if _, err := os.Stat(file); err != nil {
		file, err = expendHome(".gitconfig")
		if err != nil {
			return "", err
		}
	}
	return file, nil
}

// unsetHome unsets HOME related environments
func unsetHome() {
	if runtime.GOOS == "windows" {
		os.Unsetenv("USERPROFILE")
		os.Unsetenv("HOMEDRIVE")
		os.Unsetenv("HOMEPATH")
	}
	os.Unsetenv("HOME")
}

// setHome sets proper HOME environments
func setHome(home string) {
	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", home)
		if strings.Contains(home, ":\\") {
			slices := strings.SplitN(home, ":\\", 2)
			if len(slices) == 2 {
				os.Setenv("HOMEDRIVE", slices[0]+":")
				os.Setenv("HOMEPATH", "\\"+slices[1])
			}
		}
	} else {
		os.Setenv("HOME", home)
	}
}
