package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
func GetTag(tag string) (oid string, err error) {
	oid, err = getRef(tag, TAG_DIR)
	return oid, err
}

/*
SetHead write a reference to the HEAD file
*/
func SetHead(ref string) error {
	err := ioutil.WriteFile(HEAD_PATH, []byte(ref), 0777)
	return err
}

/*
GetHead return a reference the HEAD file
*/
func GetHead() (ref string, err error) {
	d, err := ioutil.ReadFile(HEAD_PATH)
	ref = string(d)
	return ref, err
}

/*
GetBranch returns the oid pointed by the refs/heads/ref file
*/
func GetBranch(ref string) (oid string, err error) {
	oid, err = getRef(ref, BRANCH_DIR)
	return oid, err
}

/*
SetBranch writes the oid the refs/heads/ref file should point to
*/
func SetBranch(branch string, oid string) error {
	return setRef(branch, BRANCH_DIR, oid, true, false)
}

/*
ListHeads returns the list of branches in the refs/heads directory
*/
func ListHeads() (branches []string) {
	err := filepath.Walk(BRANCH_DIR,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			b := strings.Replace(path, BRANCH_DIR+"/", "", -1)
			branches = append(branches, b)
			return nil
		})

	if err != nil {
		log.Println(err)
	}
	return branches
}
