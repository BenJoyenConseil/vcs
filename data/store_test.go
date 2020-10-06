package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestUInit(t *testing.T) {
	// given
	dir := "/tmp/"
	// when
	UInit(dir)

	// then
	if _, err := os.Stat("/tmp/.ugit"); os.IsNotExist(err) {
		t.Error(err)
	}
}

func TestHashObject(t *testing.T) {
	// given

	// when
	result := HashObbject([]byte("Hello World"))

	// then
	if string(result) != "0a4d55a8d778e5022fab701977c5d840bbc486d0" {
		t.Errorf("Fail result = %s", result)
	}
}

func TestPutObject(t *testing.T) {
	// given
	contentToVersion := "Hello World"
	expected := fmt.Sprintf("%s%s%s", "blob", string('\x00'), "Hello World")

	// when
	PutObject(contentToVersion, BLOB)

	// then
	result, _ := ioutil.ReadFile(".ugit/objects/0a6649a0077da1bf5a8b3b5dd3ea733ea6a81938")
	if bytes.Compare([]byte(expected), result) != 0 {
		t.Error("Fail. actual : ", result, " expected : ", []byte(expected))
	}

	// teardown
	os.Remove("tmp")
	os.RemoveAll(".ugit/")
}

func TestGetObject(t *testing.T) {
	// given
	oid := "0a4d55a8d778e5022fab701977c5d840bbc486d0"
	objectContent := "Hello World"
	os.MkdirAll(OBJECTS_DIR, 0777)
	ioutil.WriteFile(".ugit/objects/0a4d55a8d778e5022fab701977c5d840bbc486d0", []byte(objectContent), 0777)

	//when
	result, _type := GetObject(oid)

	// then
	if result != objectContent {
		t.Error("Data content is not correcte. actual : ", result, " != expected : ", objectContent)
	}
	if _type != "blob" {
		t.Error("Type is not correct. actual : ", _type, "expected : ", "blob")
	}

	// teardown
	os.RemoveAll(".ugit")
}
