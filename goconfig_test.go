package goconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDanyel(t *testing.T) {
	filename := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Reading file %v failed", filename)
	}
	config, lineno, err := Parse(bytes, filename)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, int(lineno))
	_ = config
	assert.Equal(t, "Danyel Bayraktar", config.Get("user.name"))
	assert.Equal(t, "cydrop@gmail.com", config.Get("user.email"))
	assert.Equal(t, "subl -w", config.Get("core.editor"))
	assert.Equal(t, `!git config --get-regexp 'alias.*' | colrm 1 6 | sed 's/[ ]/ = /' | sort`, config.Get("alias.aliases"))
}

func TestInvalidKey(t *testing.T) {
	invalidConfig := ".name = Danyel"
	config, lineno, err := Parse([]byte(invalidConfig), "")
	assert.Equal(t, ErrInvalidKeyChar, err)
	assert.Equal(t, 1, int(lineno))
	assert.Equal(t, NewGitConfig(), config)
}

func TestNoNewLine(t *testing.T) {
	validConfig := "[user] name = Danyel"
	config, lineno, err := Parse([]byte(validConfig), "")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, int(lineno))
	expect := NewGitConfig()
	expect.Add("user.name", "Danyel")
	assert.Equal(t, expect, config)
}

func TestUpperCaseKey(t *testing.T) {
	validConfig := "[core]\nQuotePath = false\n"
	config, lineno, err := Parse([]byte(validConfig), "")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, int(lineno))
	expect := NewGitConfig()
	expect.Add("core.quotepath", "false")
	assert.Equal(t, expect, config)
}

func TestExtended(t *testing.T) {
	validConfig := `[http "https://my-website.com"] sslVerify = false`
	config, lineno, err := Parse([]byte(validConfig), "")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, int(lineno))
	expect := NewGitConfig()
	expect.Add("http.https://my-website.com.sslverify", "false")
	assert.Equal(t, expect, config)
}

func ExampleParse() {
	gitconfig := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(gitconfig)
	if err != nil {
		log.Fatalf("Couldn't read file %v\n", gitconfig)
	}

	config, lineno, err := Parse(bytes, gitconfig)
	if err != nil {
		log.Fatalf("Error on line %d: %v\n", lineno, err)
	}
	fmt.Println()
	fmt.Println(lineno)
	fmt.Println(config.Get("user.name"))
	fmt.Println(config.Get("user.email"))
	// Output:
	// 10
	// Danyel Bayraktar
	// cydrop@gmail.com
}

func BenchmarkParse(b *testing.B) {
	gitconfig := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(gitconfig)
	if err != nil {
		b.Fatalf("Couldn't read file %v: %s\n", gitconfig, err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(bytes, gitconfig)
	}
}
