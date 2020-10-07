package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUInit(t *testing.T) {
	// given
	dir := "/tmp/"
	// when
	UInit(dir)

	// then
	assert.DirExists(t, "/tmp/.ugit")
}

func TestHashObject(t *testing.T) {
	// given

	// when
	result := HashObbject([]byte("Hello World"))

	// then
	assert.Equal(t, "0a4d55a8d778e5022fab701977c5d840bbc486d0", string(result))
}

func TestPutObject(t *testing.T) {
	// given
	contentToVersion := "Hello World"
	expected := fmt.Sprintf("%s%s%s", "blob", string('\x00'), "Hello World")

	// when
	PutObject(contentToVersion)

	// then
	assert.FileExists(t, ".ugit/objects/0a6649a0077da1bf5a8b3b5dd3ea733ea6a81938")
	result, _ := ioutil.ReadFile(".ugit/objects/0a6649a0077da1bf5a8b3b5dd3ea733ea6a81938")
	assert.Equal(t, result, []byte(expected))

	// teardown
	os.Remove("tmp")
	os.RemoveAll(".ugit/")
}

func TestGetObject(t *testing.T) {
	// given
	oid := "0a4d55a8d778e5022fab701977c5d840bbc486d0"
	objectContent := "Hello World"
	os.MkdirAll(OBJECTS_DIR, 0777)
	ioutil.WriteFile(".ugit/objects/"+oid, []byte("blob"+string('\x00')+objectContent), 0777)

	//when
	result, _type, err := GetObject(oid)

	// then
	assert.Nil(t, err)
	assert.Equal(t, objectContent, result)
	assert.Equal(t, _type, BLOB)

	// teardown
	os.RemoveAll(".ugit")
}
