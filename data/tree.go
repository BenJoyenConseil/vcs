package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func WriteTree(directory string) (oid string, err error) {
	tree := ""
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(directory); err != nil {
		return "", err
	}

	for _, f := range files {
		if f.Name() == ".ugit" {
			break
		}
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
