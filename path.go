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

// Abs returns absolute path and will expend homedir if path has "~/' prefix
func Abs(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return expendHome(name)
	}

	return filepath.Abs(name)
}

// AbsJoin returns absolute path, and use <dir> as parent dir for relative path
func AbsJoin(dir, name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return expendHome(name)
	}

	return Abs(filepath.Join(dir, name))
}

// IsGitDir test whether dir is a valid git dir
func IsGitDir(dir string) bool {
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

// FindGitDir searches git dir
func FindGitDir(dir string) (string, error) {
	var err error

	dir, err = Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		// Check if is in a bare repo
		if IsGitDir(dir) {
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
			if IsGitDir(gitdir) {
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
					realgit, err = AbsJoin(filepath.Dir(gitdir), realgit)
					if err != nil {
						return "", err
					}
				}
				if IsGitDir(realgit) {
					return realgit, nil
				}
				return "", fmt.Errorf("gitdir '%s' points to corrupt git repo: %s", gitdir, realgit)
			}
			return "", fmt.Errorf("bad gitdir file '%s'", gitdir)
		}
	}
	return "", nil
}

// FindGitConfig returns local git config file
func FindGitConfig(dir string) (string, error) {
	dir, err := FindGitDir(dir)
	if err == nil && dir != "" {
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
