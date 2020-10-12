package storage

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ObjectType string

const (
	UGIT_DIR       string = ".ugit"
	OBJECTS_DIR    string = ".ugit/objects"
	HEAD_PATH      string = ".ugit/HEAD"
	BYTE_SEPARATOR byte   = '\x00'
)

func UInit(dir string) {
	initDir, _ := filepath.Abs(fmt.Sprintf("%s/%s", dir, UGIT_DIR))
	log.Printf("Initialize uGit in the following directory : %s", initDir)
	err := os.MkdirAll(initDir, 0777)
	if err != nil {
		log.Println(err)
	}
}

/*
HashObject return the SHA1 result
*/
func HashObject(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	oid := []byte(fmt.Sprintf("%x", h.Sum(nil)))
	return oid
}

/*
PutObject stores the data into the ugit objects repository. An object type is added before the content inside the file.
We hash the whole to return it as the oid
*/
func PutObject(data string, objectType ...ObjectType) (oid string, err error) {
	_type := BLOB
	if len(objectType) > 0 {
		_type = objectType[0]
	}
	encoded := []byte(string(_type) + string(BYTE_SEPARATOR) + data)
	oid = string(HashObject(encoded))
	objectPath := fmt.Sprintf("%s/%s", OBJECTS_DIR, oid)
	os.MkdirAll(OBJECTS_DIR, 0777)
	err = ioutil.WriteFile(objectPath, encoded, 0777)
	return oid, err
}

/*
GetObject returns the content of the file, and its type
*/
func GetObject(oid string) (string, ObjectType, error) {
	objectPath := fmt.Sprintf("%s/%s", OBJECTS_DIR, oid)
	data, err := ioutil.ReadFile(objectPath)
	if err != nil {
		return "", ObjectType(""), err
	}
	parts := bytes.Split(data, []byte{BYTE_SEPARATOR})
	_type := ObjectType(parts[0])
	content := string(parts[1])
	return content, _type, err
}