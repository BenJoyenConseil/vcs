package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"vcs/storage"
	"vcs/tree"
)

var Commands map[string]*Command

type Command struct {
	Function      func(args []string)
	flagset       *flag.FlagSet
	parameters    map[string]string
	parsedOutputs map[string]*string
	usage         string
}

func init() {
	Commands = map[string]*Command{
		"commit": {
			Function:      commit,
			parameters:    map[string]string{"m": "The message of the commit inside double quotes"},
			parsedOutputs: map[string]*string{},
			usage:         "Snapshot the current directory in VCS and return the oid",
		},
		"init": {
			Function: uinit,
			usage:    "<directory>",
		},
		"log": {
			Function: ulog,
			usage:    "<ref>",
		},
		"checkout": {usage: "<ref>"},
	}

	for name, cmd := range Commands {
		cmd.flagset = flag.NewFlagSet(name, flag.ExitOnError)
		for p, u := range cmd.parameters {
			cmd.parsedOutputs[p] = cmd.flagset.String(p, "", u)
		}
	}

	flag.Usage = func() {
		fmt.Printf("Usage of vcs \n")
		for name, cmd := range Commands {
			fmt.Printf(" %s : %s\n", name, cmd.usage)
			cmd.flagset.PrintDefaults()
		}
	}
}

func commit(args []string) {
	Commands["commit"].flagset.Parse(args)
	if val, _ := Commands["commit"].parsedOutputs["m"]; *val != "" {
		log.Println("Commit hash : ", tree.Commit("./", *Commands["commit"].parsedOutputs["m"]))
		os.Exit(0)
	} else {
		Commands["commit"].flagset.PrintDefaults()
		os.Exit(1)
	}
}

func uinit(args []string) {
	path := "./"
	if len(os.Args) >= 3 {
		path = os.Args[2]
	}
	storage.UInit(path)
	os.Exit(0)
}

func ulog(args []string) {
	tree.PrintLog(tree.Log())
	os.Exit(0)
}
