package cli

import (
	"log"
	"os"
	"strings"
	"vcs/storage"
	"vcs/tree"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("vcs", "Snapshot your working directory")
	commit       = app.Command("commit", "snapshot the current directory with an explicite message desciption")
	commitMsg    = commit.Flag("message", "message description").Short('m').Required().String()
	vcsDir       = commit.Flag("dir", "Force to use a specific vcs directory").Default(".").Short('d').String()
	commitAction = commit.Action(ucommit)

	_init      = app.Command("init", "Setup the directory you want to be managed")
	initDir    = _init.Arg("dir", "The directory to setup VCS").Default(".").String()
	initAction = _init.Action(uinit)

	_log      = app.Command("log", "Print the commit log history")
	logaction = _log.Action(ulog)
	logRef    = _log.Arg("ref", "A reference as the starting point in the history").Default("HEAD").String()

	hash_object = app.Command("hash_object", "Save an object in vcs and get its hash")
	hashString  = hash_object.Arg("content", "The string/binary content of an object (e.g file)").Required().String()
	hashAction  = hash_object.Action(uhashObject)
	putAction   = hash_object.Flag("save", "Save into the vcs directory").Short('s').Action(uputObject).Bool()

	checkout       = app.Command("checkout", "Restore files and folders from a committed snapshot")
	oidCheckout    = checkout.Arg("oid", "The commit oid").Required().String()
	checkoutAction = checkout.Action(ucheckout)
	checkoutDir    = checkout.Flag("dir", "Force to use a specific vcs directory").Default(".").Short('d').String()

	tag       = app.Command("tag", "Give a name to the current commit")
	tagAction = tag.Action(utag)
	tagName   = tag.Arg("name", "Usefull name to find a commit oid").Required().String()

	branch       = app.Command("branch", "Print all created branches")
	branchAction = branch.Action(ubranch)
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

func uhashObject(c *kingpin.ParseContext) error {
	if !*putAction {
		data := []byte(*hashString)
		log.Println(string(storage.HashObject(data)))
	}
	return nil
}

func uputObject(c *kingpin.ParseContext) error {
	oid, err := storage.PutObject(*hashString)
	log.Println(oid)
	return err
}

func ucheckout(c *kingpin.ParseContext) error {
	err := tree.Checkout(*oidCheckout, *checkoutDir)
	return err
}

func utag(c *kingpin.ParseContext) error {
	ref, err := storage.GetHead()
	if err != nil {
		return nil
	}
	oid := tree.GetOid(ref)
	err = storage.SetTag(*tagName, oid)
	return err
}

func ubranch(c *kingpin.ParseContext) error {
	ref, err := storage.GetHead()
	if err != nil {
		return err
	}

	branches := storage.ListHeads()
	for _, b := range branches {
		if strings.Contains(ref, b) {
			log.Println("*", b)
		} else {
			log.Println(b)
		}
	}
	return nil
}
