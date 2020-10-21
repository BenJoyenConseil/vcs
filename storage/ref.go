package storage

import "io/ioutil"

/*
SetHead write the oid reference into the HEAD file
*/
func SetHead(oid string) error {
	err := ioutil.WriteFile(HEAD_PATH, []byte(oid), 0777)
	return err
}

/*
GetHead return the oid stored by the HEAD file
*/
func GetHead() (oid string) {
	d, _ := ioutil.ReadFile(HEAD_PATH)
	return string(d)
}
