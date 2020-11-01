package storage

import (
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

/*
MapOidRefs returns a index of oid and its associated references
*/
func MapOidRefs() (oidRefs map[string][]string) {
	searchPath := UGIT_DIR + "/" + REF_DIR
	oidRefs = map[string][]string{}
	headRef, _ := GetHead()
	headOid, _ := GetRef(headRef)
	oidRefs[headOid] = append([]string{}, "HEAD")

	err := filepath.Walk(
		searchPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			log.Println(path)
			rpath, _ := filepath.Rel(UGIT_DIR, path)
			oid, err := GetRef(rpath)
			if err != nil {
				return err
			}
			if strings.Contains(path, TAG_DIR) {
				rpath, _ = filepath.Rel(UGIT_DIR+"/"+TAG_DIR, path)
			} else if strings.Contains(path, BRANCH_DIR) {
				rpath, _ = filepath.Rel(UGIT_DIR+"/"+BRANCH_DIR, path)
			}
			oidRefs[oid] = append(oidRefs[oid], rpath)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return oidRefs
}
