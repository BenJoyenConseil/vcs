package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func WriteTree(directory string) (string, error) {
	tree := ""
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println(err)
		return tree, err
	}
	for _, f := range files {
		if f.Name() == ".ugit" {
			break
		}
		path := fmt.Sprintf("%s/%s", directory, f.Name())
		if f.IsDir() {
			oid, _ := WriteTree(path)
			tree += fmt.Sprintf("%s %s %s\n", "tree", oid, f.Name())
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
			}
			oid := PutObject(string(data), BLOB)
			tree += fmt.Sprintf("%s %s %s\n", "blob", oid, f.Name())
		}
	}
	oid := PutObject(tree, TREE)
	log.Println(oid)
	return oid, nil
}

func ReadTree(oid string, basePath ...string) error {
	data, _type, err := GetObject(oid)
	path := "."
	if len(basePath) > 0 {
		path = basePath[0]
	}
	if _type != TREE {
		return nil
	}
	for _, line := range strings.Split(data, "\n") {
		lineSplits := strings.Split(line, " ")

		if ObjectType(lineSplits[0]) == TREE {
			dir := fmt.Sprintf("%s/%s", path, lineSplits[2])
			os.Mkdir(dir, 0777)
			ReadTree(lineSplits[1], dir)
		} else {
			d, _, err := GetObject(lineSplits[1])
			if err != nil {
				return err
			}
			ioutil.WriteFile(fmt.Sprintf("%s/%s", path, lineSplits[2]), []byte(d), 0777)
		}
	}
	return err
}
