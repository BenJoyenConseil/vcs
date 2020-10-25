package tree

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"vcs/storage"
)

var IGNORED_PATH = []string{
	".ugit",
	".git",
	".gitignore",
	"./",
}

/*
CommitNode represent a commit's linkedlist
*/
type CommitNode struct {
	oid     string
	parent  *CommitNode
	message string
}

/*
IsIgnored return true if the path should be ignored by the VCS
*/
func IsIgnored(path string) bool {
	for _, ignore := range IGNORED_PATH {
		if strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}

/*
WriteTree saves the directory in the content object storage with recursion
*/
func WriteTree(directory string) (oid string, err error) {
	tree := ""
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(directory); err != nil {
		return "", err
	}
	for _, f := range files {
		if IsIgnored(f.Name()) {
			continue
		}
		log.Println(f.Name())
		path := fmt.Sprintf("%s/%s", directory, f.Name())
		if f.IsDir() {
			oid, err = WriteTree(path)
			tree += fmt.Sprintf("%s %s %s\n", storage.TREE, oid, f.Name())
		} else {
			var data []byte
			if data, err = ioutil.ReadFile(path); err != nil {
				return "", err
			}
			oid, err = storage.PutObject(string(data), storage.BLOB)
			tree += fmt.Sprintf("%s %s %s\n", storage.BLOB, oid, f.Name())
		}
		if err != nil {
			return "", err
		}
	}
	oid, err = storage.PutObject(tree, storage.TREE)
	return oid, err
}

/*
ReadTree restores the tree contents referenced by the oid into the basepath
*/
func ReadTree(oid string, basePath ...string) error {
	path := "."
	if len(basePath) > 0 {
		path = basePath[0]
	}
	log.Printf("Restoring tree {%s} in location %s", oid, path)
	data, _type, err := storage.GetObject(oid)
	if err != nil {
		return err
	}
	if _type != storage.TREE {
		return errors.New("This oid point to a non tree object : " + string(_type))
	}
	if !IsIgnored(path) {
		log.Println("Remove : ", path, os.RemoveAll(path))
	}
	os.Mkdir(path, 0777)

	treeLines := strings.Split(data, "\n")
	for _, line := range treeLines {
		lineSplits := strings.Split(line, " ")
		if len(lineSplits) != 3 {
			continue
		}
		t := storage.ObjectType(lineSplits[0])
		o := lineSplits[1]
		p := lineSplits[2]

		if t == storage.TREE {
			subdir := fmt.Sprintf("%s/%s", path, p)
			if err := ReadTree(o, subdir); err != nil {
				return err
			}
		} else {
			d, _, err := storage.GetObject(o)
			if err != nil {
				return err
			}
			filePath := fmt.Sprintf("%s/%s", path, p)
			log.Println("Creating file ", filePath)
			ioutil.WriteFile(filePath, []byte(d), 0777)
		}
	}
	return nil
}

/*
Commit takes a snapshot of the directory and add message plus metadata
*/
func Commit(dir string, message string, metadata ...string) (oid string) {
	log.Println("Snapshoting the following directory : ", dir)
	oid, err := WriteTree(dir)
	if err != nil {
		log.Println(err)
		return ""
	}

	ref, err := storage.GetHead()
	if err != nil {
		log.Println(err)
		return ""
	}

	oidParent := GetOid(ref)
	commit := fmt.Sprintf("%s %s\n", storage.TREE, oid)
	commit += fmt.Sprintf("%s %s\n", storage.PARENT, oidParent)
	commit += fmt.Sprintf("\n%s", message)

	oid, err = storage.PutObject(commit, storage.COMMIT)
	if err != nil {
		log.Println(err)
		return ""
	}

	storage.SetBranch(ref, oid)
	return oid
}

/*
GetCommit return the tree oid, its parent's commit oid, and its associated message
*/
func GetCommit(oid string) (tree string, parent string, message string, err error) {
	data, t, err := storage.GetObject(oid)
	if err != nil {
		return "", "", "", err
	}
	if t != storage.COMMIT {
		return "", "", "", errors.New("The object " + oid + " is not a commit : " + string(t))
	}
	commitLines := strings.Split(data, "\n")
	for i, l := range commitLines {
		token := strings.Split(l, " ")
		if i < 2 {
			_type := storage.ObjectType(token[0])
			_oid := token[1]
			if _type == storage.PARENT {
				parent = _oid
			} else if _type == storage.TREE {
				tree = _oid
			}
		} else {
			message += l
		}
	}
	return tree, parent, message, nil
}

/*
Log iterates over the commits to build a linked list structure starting from the HEAD to the first Commit
*/
func Log(ref ...string) *CommitNode {
	var currentNode *CommitNode
	var headNode *CommitNode

	r, err := storage.GetHead()
	oid := GetOid(r)
	if err != nil {
		log.Println(err)
		return nil
	}

	_, parent, message, err := GetCommit(oid)
	if err != nil {
		log.Println(err)
		return nil
	}
	for oid != "" {

		previous := currentNode
		currentNode = &CommitNode{
			message: message,
			oid:     oid,
		}
		if previous != nil {
			previous.parent = currentNode
		} else {
			headNode = currentNode
		}
		oid = parent
		_, parent, message, _ = GetCommit(oid)
	}
	return headNode
}

/*
PrintLog takes a linked list of commits and prints them
*/
func PrintLog(commit *CommitNode) {
	current := commit
	for current != nil {
		char := "|"
		if current.parent == nil {
			char = " "
		}
		fmt.Println("*", "commit", current.oid)
		fmt.Println(char)
		fmt.Println(char, "\t", current.message)
		fmt.Println(char)
		current = current.parent
	}
}

/*
Checkout  moves HEAD to a commit oid and restore its state (e.g files and folders)
*/
func Checkout(ref string, basePath ...string) error {
	oid := GetOid(ref)
	storage.SetHead(ref)
	tree, _, _, err := GetCommit(oid)
	if err != nil {
		return err
	}
	return ReadTree(tree, basePath...)
}

/*
GetOid find the oid of a commit, a reference, or a tag
*/
func GetOid(ref string) (oid string) {
	if oid, err := storage.GetTag(ref); err == nil {
		return oid
	}
	if oid, err := storage.GetBranch(ref); err == nil {
		return oid
	}
	return ref
}
