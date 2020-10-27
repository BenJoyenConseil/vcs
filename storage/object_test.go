package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"vcs/mock"

	"github.com/stretchr/testify/assert"
)

func TestUInit(t *testing.T) {
	// given
	dir := "./"

	// when
	UInit(dir)

	// then
	assert.DirExists(t, ".ugit")

	assert.FileExists(t, ".ugit/HEAD")
	d, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "refs/heads/master", string(d))

	assert.DirExists(t, ".ugit/refs/heads/")
	assert.FileExists(t, ".ugit/refs/heads/master")
	d, _ = ioutil.ReadFile(".ugit/refs/heads/master")
	assert.Equal(t, "", string(d))

	assert.DirExists(t, ".ugit/objects")

	mock.Teardown()
}
func TestUInit_GracefullWhenAllreadyExists(t *testing.T) {
	// given
	mock.SetupUgitDir()

	// when
	UInit(".")

	// then
	assert.DirExists(t, ".ugit")
	assert.FileExists(t, ".ugit/HEAD")
	assert.DirExists(t, ".ugit/refs/heads/")
	assert.FileExists(t, ".ugit/refs/heads/master")
	assert.FileExists(t, ".ugit/refs/tags/v0.1.0")

	d, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "refs/heads/master", string(d))

	d, _ = ioutil.ReadFile(".ugit/refs/heads/master")
	assert.Equal(t, "cdf776713053cc0710735a61dfbe6492f3ed31b2", string(d))
	d, _ = ioutil.ReadFile(".ugit/refs/tags/v0.1.0")
	assert.Equal(t, "93584d4997160f16e3ac4390ec4008a2d2ff32d6", string(d))

	mock.Teardown()
}

func TestHashObject(t *testing.T) {
	// given

	// when
	result := HashObject([]byte("Hello World"))

	// then
	assert.Equal(t, "0a4d55a8d778e5022fab701977c5d840bbc486d0", string(result))
}

func TestPutObject(t *testing.T) {
	// given
	contentToVersion := "Hello World"
	os.MkdirAll(".ugit/objects", 0777)
	expected := fmt.Sprintf("%s%s%s", "blob", string('\x00'), "Hello World")

	// when
	PutObject(contentToVersion)

	// then
	assert.FileExists(t, ".ugit/objects/0a6649a0077da1bf5a8b3b5dd3ea733ea6a81938")
	result, _ := ioutil.ReadFile(".ugit/objects/0a6649a0077da1bf5a8b3b5dd3ea733ea6a81938")
	assert.Equal(t, result, []byte(expected))

	// teardown
	mock.Teardown()
}

func TestGetObject(t *testing.T) {
	// given
	oid := "0a4d55a8d778e5022fab701977c5d840bbc486d0"
	objectContent := "Hello World"
	os.MkdirAll(".ugit/objects/", 0777)
	ioutil.WriteFile(".ugit/objects/"+oid, []byte("blob"+string('\x00')+objectContent), 0777)

	//when
	result, _type, err := GetObject(oid)

	// then
	assert.Nil(t, err)
	assert.Equal(t, objectContent, result)
	assert.Equal(t, _type, BLOB)

	// teardown
	mock.Teardown()
}
