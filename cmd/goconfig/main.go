package main

import (
	"fmt"
	"os"

	"github.com/jiangxin/goconfig"
	flag "github.com/spf13/pflag"
)

var (
	optGlobal         bool
	optSystem         bool
	optLocal          bool
	optFilename       string
	optInclude        bool
	optActionGet      bool
	optActionGetAll   bool
	optActionAdd      bool
	optActionSet      bool
	optActionUnset    bool
	optActionUnsetAll bool
	optActionList     bool

	configFile string
	cfg        goconfig.GitConfig
)

func checkOptions() error {
	var (
		scopes      = 0
		actions     = 0
		writeAction bool
		err         error
	)

	if optSystem {
		configFile = goconfig.SystemConfigFile()
		scopes++
	}
	if optGlobal {
		configFile, err = goconfig.GlobalConfigFile()
		if err != nil {
			return err
		}
		scopes++
	}
	if optLocal {
		configFile, err = goconfig.FindGitConfig("")
		if err != nil {
			return err
		}
		scopes++
	}
	if optFilename != "" {
		configFile = optFilename
		scopes++
	}
	if scopes > 1 {
		return fmt.Errorf("only one config file at a time")
	}

	if optActionGet {
		actions++
	}
	if optActionGetAll {
		actions++
	}
	if optActionList {
		actions++
	}
	if optActionAdd {
		writeAction = true
		actions++
	}
	if optActionUnset {
		writeAction = true
		actions++
	}
	if optActionUnsetAll {
		writeAction = true
		actions++
	}
	if actions == 0 {
		if len(flag.Args()) == 1 {
			optActionGet = true
		} else if len(flag.Args()) == 2 {
			writeAction = true
			optActionSet = true
		} else {
			return fmt.Errorf("wrong number of arguments, should be 1 or 2")
		}
	}
	if actions > 1 {
		return fmt.Errorf("only one action at a time")
	}
	if configFile == "" {
		configFile, err = goconfig.FindGitConfig("")
		if err != nil {
			if err != goconfig.ErrNotInGitDir || writeAction {
				return err
			}
		}
		optInclude = true
	}

	return nil
}

func runGet(args ...string) error {
	for _, k := range args {
		fmt.Println(cfg.Get(k))
	}
	return nil
}
func runGetAll(args ...string) error {
	for _, k := range args {
		for _, v := range cfg.GetAll(k) {
			fmt.Println(v)
		}
	}
	return nil
}

func runList(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("wrong number of arguments, should be 0")
	}
	for _, k := range cfg.Keys() {
		for _, v := range cfg.GetAll(k) {
			fmt.Printf("%s=%s\n", k, v)
		}
	}
	return nil
}

func runAdd(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments, should be 2")
	}
	cfg.Add(args[0], args[1])
	return cfg.Save(configFile)
}

func runSet(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments, should be 2")
	}
	cfg.Set(args[0], args[1])
	return cfg.Save(configFile)
}

func runUnset(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong number of arguments, should be 1")
	}
	cfg.Unset(args[0])
	return cfg.Save(configFile)
}

func runUnsetAll(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong number of arguments, should be 1")
	}
	cfg.UnsetAll(args[0])
	return cfg.Save(configFile)
}

func main() {
	var err error

	err = checkOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	if optInclude {
		cfg, err = goconfig.LoadAll(configFile)
	} else {
		cfg, err = goconfig.Load(configFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	if optActionGet {
		err = runGet(flag.Args()...)
	} else if optActionGetAll {
		err = runGetAll(flag.Args()...)
	} else if optActionList {
		err = runList(flag.Args()...)
	} else if optActionAdd {
		err = runAdd(flag.Args()...)
	} else if optActionSet {
		err = runSet(flag.Args()...)
	} else if optActionUnset {
		err = runUnset(flag.Args()...)
	} else if optActionUnsetAll {
		err = runUnsetAll(flag.Args()...)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	// scope option
	flag.BoolVar(&optGlobal, "global", false, "use global config file")
	flag.BoolVar(&optSystem, "system", false, "use system config file")
	flag.BoolVar(&optLocal, "local", false, "use local config file")
	flag.BoolVar(&optInclude, "include", false, "respect include directives on lookup")
	flag.StringVarP(&optFilename, "file", "f", "", "file to load")
	// action option
	flag.BoolVar(&optActionGet, "get", false, "get value: name")
	flag.BoolVar(&optActionGetAll, "get-all", false, "get value: name")
	flag.BoolVar(&optActionAdd, "add", false, "gdd a new variable: name value")
	flag.BoolVar(&optActionUnset, "unset", false, "remove a variable")
	flag.BoolVar(&optActionUnsetAll, "unset-all", false, "remove all matches")
	flag.BoolVarP(&optActionList, "list", "l", false, "list all")
	flag.Parse()
}
