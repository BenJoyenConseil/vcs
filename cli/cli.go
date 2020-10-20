package cli

import (
	"os"
	"vcs/storage"
	"vcs/tree"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("vcs", "Snapshot your working directory")
	commit       = app.Command("commit", "snapshot the current directory with an explicite message desciption")
	commitMsg    = commit.Flag("message", "message description").Short('m').Required().String()
	vcsDir       = commit.Flag("dir", "Force to use a specific vcs directory").Default("./").Short('d').String()
	commitAction = commit.Action(ucommit)

	_init      = app.Command("init", "Setup the directory you want to be managed")
	initDir    = _init.Arg("dir", "The directory to setup VCS").Default("./").String()
	initAction = _init.Action(uinit)

	_log      = app.Command("log", "Print the commit log history")
	logaction = _log.Action(ulog)
)

func Parse(args []string) {
	kingpin.MustParse(app.Parse(os.Args[1:]))

}

func ucommit(c *kingpin.ParseContext) error {
	tree.Commit(*vcsDir, *commitMsg)
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
