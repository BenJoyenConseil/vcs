package cli

import (
	"fmt"
	"os"
	"vcs/storage"
	"vcs/tree"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("vcs", "Version your directory. period")
	commit       = app.Command("commit", "snapshot the current directory")
	commitMsg    = commit.Flag("message", "message description").Short('m').Required().String()
	commitAction = commit.Action(ucommit)

	_init      = app.Command("init", "Setup directory to be manageed")
	initDir    = _init.Arg("dir", "The directory to setup VCS").Default("./").String()
	initAction = _init.Action(uinit)

	_log      = app.Command("log", "Print the commit log history")
	logaction = _log.Action(ulog)
)

func Init(args []string) {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	fmt.Println(*commitMsg, *initDir)

}

func ucommit(c *kingpin.ParseContext) error {
	tree.Commit("./", *commitMsg)
	return nil
}

func uinit(c *kingpin.ParseContext) error {
	storage.UInit(*initDir)
	return nil
}

func ulog(c *kingpin.ParseContext) error {
	tree.PrintLog(tree.Log())
	return nil
}
