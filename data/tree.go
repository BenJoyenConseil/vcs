package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var IGNORED_PATH = map[string]struct{}{
	".ugit":      struct{}{},
	".git":       struct{}{},
	".gitignore": struct{}{},
}

type CommitNode struct {
	oid     string
	parent  *CommitNode
	message string
}

func IsIgnored(path string) bool {
	if _, ok := IGNORED_PATH[path]; ok {
		return true
	}
	return false
}

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
			tree += fmt.Sprintf("%s %s %s\n", TREE, oid, f.Name())
		} else {
			var data []byte
			if data, err = ioutil.ReadFile(path); err != nil {
				return "", err
			}
			oid, err = PutObject(string(data), BLOB)
			tree += fmt.Sprintf("%s %s %s\n", BLOB, oid, f.Name())
		}
		if err != nil {
			return "", err
		}
	}
	oid, err = PutObject(tree, TREE)
	return oid, err
}

func ReadTree(oid string, basePath ...string) error {
	path := "."
	if len(basePath) > 0 {
		path = basePath[0]
	}
	log.Printf("Restoring tree {%s} in location %s", oid, path)
	data, _, err := GetObject(oid)
	if err != nil {
		return err
	}
	os.RemoveAll(path)
	os.Mkdir(path, 0777)

	treeLines := strings.Split(data, "\n")
	for _, line := range treeLines {
		lineSplits := strings.Split(line, " ")
		t := ObjectType(lineSplits[0])
		o := lineSplits[1]
		p := lineSplits[2]

		if t == TREE {
			subdir := fmt.Sprintf("%s/%s", path, p)
			if err := ReadTree(o, subdir); err != nil {
				return err
			}
		} else {
			d, _, err := GetObject(o)
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

func Commit(dir string, message string, metadata ...string) (oid string) {
	log.Println(dir)
	oid, err := WriteTree(dir)
	if err != nil {
		log.Println(err)
		return ""
	}
	commit := fmt.Sprintf("%s %s\n", TREE, oid)
	commit += fmt.Sprintf("%s %s\n", PARENT, GetHead())
	commit += fmt.Sprintf("\n%s", message)

	oid, err = PutObject(commit, COMMIT)
	if err != nil {
		log.Println(err)
		return ""
	}
	SetHead(oid)
	return oid
}

func GetCommit(oid string) (tree string, parent string, message string, err error) {
	data, t, err := GetObject(oid)
	if err != nil {
		return "", "", "", err
	}
	if t != COMMIT {
		return "", "", "", errors.New("The object " + oid + " is not a commit : " + string(t))
	}
	commitLines := strings.Split(data, "\n")
	for i, l := range commitLines {
		token := strings.Split(l, " ")
		if i < 2 {
			_type := ObjectType(token[0])
			_oid := token[1]
			if _type == PARENT {
				parent = _oid
			} else if _type == TREE {
				tree = _oid
			}
		} else {
			message += l
		}
	}
	return tree, parent, message, nil
}

func SetHead(oid string) error {
	err := ioutil.WriteFile(HEAD, []byte(oid), 0777)
	return err
}

var HEAD = fmt.Sprintf("%s/HEAD", UGIT_DIR)

func GetHead() (oid string) {
	d, _ := ioutil.ReadFile(HEAD)
	return string(d)
}

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
