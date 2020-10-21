package storage

import (
	"errors"
	"io/ioutil"
	"os"
)

/*
SetRef write the oid reference into the reference's file
*/
func SetRef(ref string, oid string) error {
	os.MkdirAll(REF_DIR, 0777)
	err := ioutil.WriteFile(ref, []byte(oid), 0777)
	return err
}

/*
GetRef return the oid stored by the reference's file
*/
func GetRef(ref string) (oid string, err error) {
	d, err := ioutil.ReadFile(ref)
	return string(d), err
}

/*
SetTag write the oid reference into the tag's file
*/
func SetTag(tag string, oid string) error {
	os.MkdirAll(TAG_DIR, 0777)
	tag_path := TAG_DIR + "/" + tag
	if _, err := os.Stat(tag_path); err == nil {
		return errors.New("The tag already exists")
	}
	err := SetRef(tag_path, oid)
	return err
}

/*
GetTag return the oid stored by the tag's file
*/
func GetTag(tag string) (oid string) {
	oid, _ = GetRef(TAG_DIR + "/" + tag)
	return oid
}

/*
SetHead write the oid reference into the HEAD file
*/
func SetHead(oid string) {
	SetRef(HEAD_PATH, oid)
}

/*
GetHead return the oid stored by the HEAD file
*/
func GetHead() (oid string) {
	oid, _ = GetRef(HEAD_PATH)
	return oid
}
