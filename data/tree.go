package data

import (
	"fmt"
	"io/ioutil"
	"log"
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
