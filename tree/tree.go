package tree

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"ugit/storage"
)

var IGNORED_PATH = map[string]struct{}{
	".ugit":      struct{}{},
	".git":       struct{}{},
	".gitignore": struct{}{},
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
	if _, ok := IGNORED_PATH[path]; ok {
		return true
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
	data, _, err := storage.GetObject(oid)
	if err != nil {
		return err
	}
	os.RemoveAll(path)
	os.Mkdir(path, 0777)

	treeLines := strings.Split(data, "\n")
	for _, line := range treeLines {
		lineSplits := strings.Split(line, " ")
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
	log.Println(dir)
	oid, err := WriteTree(dir)
	if err != nil {
		log.Println(err)
		return ""
	}
	commit := fmt.Sprintf("%s %s\n", storage.TREE, oid)
	commit += fmt.Sprintf("%s %s\n", storage.PARENT, GetHead())
	commit += fmt.Sprintf("\n%s", message)

	oid, err = storage.PutObject(commit, storage.COMMIT)
	if err != nil {
		log.Println(err)
		return ""
	}
	SetHead(oid)
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
SetHead write the oid reference into the HEAD file
*/
func SetHead(oid string) error {
	err := ioutil.WriteFile(storage.HEAD_PATH, []byte(oid), 0777)
	return err
}

/*
GetHead return the oid stored by the HEAD file
*/
func GetHead() (oid string) {
	d, _ := ioutil.ReadFile(storage.HEAD_PATH)
	return string(d)
}

/*
Log iterates over the commits to build a linked list structure starting from the HEAD to the first Commit
*/
func Log() *CommitNode {
	var currentNode *CommitNode
	var headNode *CommitNode

	oid := GetHead()
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
		log.Println("*", "commit", current.oid)
		log.Println(char)
		log.Println(char, "\t", current.message)
		log.Println(char)
		current = current.parent
	}
}
