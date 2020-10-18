package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"vcs/storage"
	"vcs/tree"
)

var commands = map[string]*flag.FlagSet{
	"commit": flag.NewFlagSet("commit", flag.ExitOnError),
	"init":   flag.NewFlagSet("init", flag.ExitOnError),
	"log":    flag.NewFlagSet("log", flag.ExitOnError),
}

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	for k, v := range commands {
		fmt.Fprintf(os.Stderr, "%s\n", k)
		v.PrintDefaults()
	}
	flag.PrintDefaults()
}

var commitMessage *string

func init() {
	commitMessage = commands["commit"].String("m", "", "The message of the commit inside double quotes")
}

func main() {
	flag.Usage = usage
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if _, ok := commands[os.Args[1]]; !ok {
		flag.Usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "commit":
		commands["commit"].Parse(os.Args[2:])
		if *commitMessage != "" {
			log.Println("Commit hash : ", tree.Commit("./", *commitMessage))
			return
		}
	case "init":
		path := "./"
		if len(os.Args) >= 3 {
			path = os.Args[2]
		}
		storage.UInit(path)
		return
	case "log":
		tree.PrintLog(tree.Log())
		return
	}

	flag.Usage()
}
