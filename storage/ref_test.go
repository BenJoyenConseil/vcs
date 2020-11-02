package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"vcs/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetRef(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/tags/", 0777)

	// when
	SetRef("refs/tags/v0.1.0", "hashcommit123")

	//then
	d, _ := ioutil.ReadFile(".ugit/refs/tags/v0.1.0")
	assert.FileExists(t, ".ugit/refs/tags/v0.1.0")
	assert.Equal(t, "hashcommit123", string(d))
	mock.Teardown()
}
func TestSetHead(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)

	// when
	SetHead("refs/heads/master")

	// then
	assert.FileExists(t, ".ugit/HEAD")
	h, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "refs/heads/master", string(h))
	mock.Teardown()
}

func TestGetHead(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("refs/heads/master"), 0777)

	// when
	oid, err := GetHead()
	os.Remove(".ugit/HEAD")
	_, err2 := GetHead()

	// then
	assert.Equal(t, "refs/heads/master", oid)
	assert.Nil(t, err)
	assert.NotNil(t, err2)
	mock.Teardown()
}

func TestListHeads(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/heads/feature", 0777)
	ioutil.WriteFile(".ugit/refs/heads/master", []byte(""), 0777)
	ioutil.WriteFile(".ugit/refs/heads/feature/yolo", []byte(""), 0777)

	// when
	branches := ListHeads()

	// then
	assert.ElementsMatch(t, branches, []string{
		"master",
		"feature/yolo",
	})
	mock.Teardown()
}

func TestListRefs(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/heads/feature", 0777)
	os.MkdirAll(".ugit/refs/tags", 0777)
	ioutil.WriteFile(".ugit/refs/heads/master", []byte("123"), 0777)
	ioutil.WriteFile(".ugit/refs/heads/feature/yolo", []byte("124"), 0777)
	ioutil.WriteFile(".ugit/refs/tags/v0.1.0", []byte("123"), 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("refs/heads/feature/yolo"), 0777)

	// when
	refs := MapOidRefs()

	// then
	assert.Equal(t, []string{"master", "v0.1.0"}, refs["123"])
	assert.Equal(t, []string{"HEAD", "feature/yolo"}, refs["124"])

	// when HEAD detached
	ioutil.WriteFile(".ugit/HEAD", []byte("123"), 0777)
	refs = MapOidRefs()
	assert.Equal(t, []string{"HEAD", "master", "v0.1.0"}, refs["123"])

	mock.Teardown()
}
