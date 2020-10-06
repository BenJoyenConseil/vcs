package data

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ObjectType string

const (
	UGIT_DIR    string     = ".ugit"
	OBJECTS_DIR string     = ".ugit/objects"
	BLOB        ObjectType = "blob"
	TREE        ObjectType = "tree"
)

func UInit(dir string) {
	initDir := filepath.Clean(fmt.Sprintf("%s/%s", dir, UGIT_DIR))
	log.Printf("Initialize uGit in the following directory : %s", initDir)
	err := os.MkdirAll(initDir, 0777)
	if err != nil {
		log.Println(err)
	}
}

func HashObbject(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	oid := []byte(fmt.Sprintf("%x", h.Sum(nil)))
	return oid
}

func PutObject(data string, objectType ObjectType) (oid string) {
	if objectType == "" {
		objectType = BLOB
	}
	encoded := []byte(string(objectType) + string('\x00') + data)
	oid = string(HashObbject(encoded))
	objectPath := fmt.Sprintf("%s/%s", OBJECTS_DIR, oid)
	os.MkdirAll(OBJECTS_DIR, 0777)
	err := ioutil.WriteFile(objectPath, encoded, 0777)
	if err != nil {
		log.Println(err)
	}
	return oid
}

func GetObject(oid string) (string, error) {
	objectPath := fmt.Sprintf("%s/%s", OBJECTS_DIR, oid)
	data, err := ioutil.ReadFile(objectPath)
	return string(data), err
}
