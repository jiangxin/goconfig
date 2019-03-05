package goconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
