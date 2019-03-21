goconfig
========

[![Travis Build Status](https://travis-ci.org/jiangxin/goconfig.svg?branch=master)](https://travis-ci.org/jiangxin/goconfig)

# Table of contents

1. Introduction
2. Usage
3. Contributing
4. Reporting bugs

-------------------

# 1. Introduction

This project parses config files that have the same syntax as gitconfig files. It understands
multiple values configuration, and can parse included configs by `include.path` directions
(`includeIf.*.path` configuration is not supported yet).

It has no knowledge of git-specific keys and as such, does not provide any convenience methods
like  `config.GetUserName()`. For these, look into [go-gitconfig](https://github.com/tcnksm/go-gitconfig)

Most of the code was copied and translated to Go from [git/config.c](https://github.com/git/git/blob/95ec6b1b3393eb6e26da40c565520a8db9796e9f/config.c)

# 2. Usage

To load specific git config file and inherit global and system git config, using:

```go
package main

import (
	"fmt"
	"log"

	"github.com/jiangxin/goconfig"
)

func main() {
	cfg, err := goconfig.LoadAll("")
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("cfg is nil")
	}

	fmt.Printf(cfg.Get("user.name"))
}
```

As an example, there is a full functional `git config` to read/write git
config file implemented by goconfig, see:

    cmd/goconfig/main.go

# 3. Contributing

Contributions are welcome! Fork -> Push -> Pull request.

# 4. Bug report / suggestions

Just create an issue! I will try to reply as soon as possible.
