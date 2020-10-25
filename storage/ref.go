package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
SetRef write the oid reference into the reference's file
*/
func setRef(ref string, refPath string, oid string, override bool, create bool) error {
	var p string
	if strings.Contains(".ugit/"+ref, refPath) {
		p = ".ugit/" + ref
	} else {
		p = fmt.Sprintf("%s/%s", refPath, ref)
	}
	log.Println(ref, refPath, p, oid, override)
	if _, notExists := os.Stat(p); notExists == nil && !override {
		return fmt.Errorf("The reference %s already exists and override is not allowed", ref)
	}
	if _, notExists := os.Stat(p); notExists != nil && !create {
		return fmt.Errorf("The reference %s does not exist and. You must create it before", ref)
	}
	os.MkdirAll(REF_DIR, 0777)
	err := ioutil.WriteFile(p, []byte(oid), 0777)
	return err
}

/*
GetRef return the oid stored by the reference's file
*/
func getRef(ref string, refPath string) (oid string, err error) {
	var p string
	if strings.Contains(".ugit/"+ref, refPath) {
		p = ".ugit/" + ref
	} else {
		p = fmt.Sprintf("%s/%s", refPath, ref)
	}
	d, err := ioutil.ReadFile(p)
	return string(d), err
}

/*
SetTag write the oid reference into the tag's file
*/
func SetTag(tag string, oid string) error {
	os.MkdirAll(TAG_DIR, 0777)
	err := setRef(tag, TAG_DIR, oid, false, true)
	return err
}

/*
GetTag return the oid stored by the tag's file
*/
func GetTag(tag string) (oid string) {
	oid, _ = getRef(tag, TAG_DIR)
	return oid
}

/*
SetHead write a reference to the HEAD file
*/
func SetHead(ref string) {
	setRef("HEAD", HEAD_PATH, ref, true, true)
}

/*
GetHead return a reference the HEAD file
*/
func GetHead() (ref string) {
	ref, _ = getRef("HEAD", HEAD_PATH)
	return ref
}

/*
GetBranch return the oid pointed by the refs/heads/ref file
*/
func GetBranch(ref string) (oid string) {
	oid, _ = getRef(ref, BRANCH_DIR)
	return oid
}

/*
SetBranch write the oid the refs/heads/ref file should point to
*/
func SetBranch(branch string, oid string) error {
	return setRef(branch, BRANCH_DIR, oid, true, false)
}
