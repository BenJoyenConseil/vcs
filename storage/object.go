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

func UInit(dir string) {
	initDir, _ := filepath.Abs(fmt.Sprintf("%s/%s", dir, UGIT_DIR))
	log.Printf("Initialize uGit in the following directory : %s", initDir)
	os.MkdirAll(initDir, 0777)

	headP, _ := filepath.Abs(fmt.Sprintf("%s/%s/HEAD", dir, UGIT_DIR))
	if _, err := os.Stat(headP); err != nil {
		ioutil.WriteFile(headP, []byte("refs/heads/master"), 0777)
	}

	branchesP, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s", dir, UGIT_DIR, BRANCH_DIR))
	os.MkdirAll(branchesP, 0777)

	tagsP, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s", dir, UGIT_DIR, TAG_DIR))
	os.MkdirAll(tagsP, 0777)

	masterP, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/master", dir, UGIT_DIR, BRANCH_DIR))
	if _, err := os.Stat(headP); err != nil {
		ioutil.WriteFile(masterP, nil, 0777)
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
	objectPath := fmt.Sprintf("%s/%s/%s", UGIT_DIR, OBJECTS_DIR, oid)
	os.MkdirAll(UGIT_DIR+"/"+OBJECTS_DIR, 0777)
	err = ioutil.WriteFile(objectPath, encoded, 0777)
	return oid, err
}

/*
GetObject returns the content of the file, and its type
*/
func GetObject(oid string) (string, ObjectType, error) {
	objectPath := fmt.Sprintf("%s/%s/%s", UGIT_DIR, OBJECTS_DIR, oid)
	data, err := ioutil.ReadFile(objectPath)
	if err != nil {
		return "", ObjectType(""), err
	}
	parts := bytes.Split(data, []byte{BYTE_SEPARATOR})
	_type := ObjectType(parts[0])
	content := string(parts[1])
	return content, _type, err
}
