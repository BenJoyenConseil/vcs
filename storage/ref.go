package storage

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
GetRef return the oid stored by the reference's file
*/
func GetRef(ref string) (oid string, err error) {
	path := UGIT_DIR + "/" + ref
	d, err := ioutil.ReadFile(path)
	return string(d), err
}

/*
SetRef write the oid reference into the reference's file
*/
func SetRef(ref string, oid string) error {
	err := ioutil.WriteFile(UGIT_DIR+"/"+ref, []byte(oid), 0777)
	return err
}

/*
SetTag write the oid reference into the tag's file
*/
func SetTag(tag string, oid string) error {
	if _, err := GetRef(tag); err != nil {
		err := SetRef(tag, oid)
		return err
	}

	return errors.New("This Tag already exists")
}

/*
GetTag return the oid stored by the tag's file
*/
func GetTag(tag string) (oid string, err error) {
	oid, err = GetRef(tag)
	return oid, err
}

/*
SetHead write a reference to the HEAD file
*/
func SetHead(ref string) error {
	err := ioutil.WriteFile(UGIT_DIR+"/"+HEAD_PATH, []byte(ref), 0777)
	return err
}

/*
GetHead return a reference the HEAD file
*/
func GetHead() (ref string, err error) {
	d, err := ioutil.ReadFile(UGIT_DIR + "/" + HEAD_PATH)
	ref = string(d)
	return ref, err
}

/*
GetBranch returns the oid pointed by the refs/heads/ref file
*/
func GetBranch(branch string) (oid string, err error) {
	oid, err = GetRef(branch)
	return oid, err
}

/*
SetBranch writes the oid the refs/heads/ref file should point to
*/
func SetBranch(branch string, oid string) error {
	if _, err := GetRef(branch); err == nil {
		return SetRef(branch, oid)
	}
	return errors.New("This branch does not exist. use 'branch' command")
}

/*
ListHeads returns the list of branches in the refs/heads directory
*/
func ListHeads() (branches []string) {
	err := filepath.Walk(UGIT_DIR+"/"+BRANCH_DIR,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			b := strings.Split(path, BRANCH_DIR+"/")[1]
			branches = append(branches, b)
			return nil
		})

	if err != nil {
		log.Println(err)
	}
	return branches
}
