package main

import (
	"flag"
	"os"
	"vcs/cli"
)

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if _, ok := cli.Commands[os.Args[1]]; !ok {
		flag.Usage()
		os.Exit(1)
	}
	cli.Commands[os.Args[1]].Function(os.Args[2:])

	flag.Usage()
}
