package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"vcs/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetHead(t *testing.T) {
	// given
	ref := "refs/heads/master"
	os.MkdirAll(".ugit/", 0777)

	// when
	SetHead(ref)

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

func TestSetTag(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)
	commitOid := "123"

	// when
	err1 := SetTag("v0.1.0", commitOid)
	err2 := SetTag("v0.1.0", commitOid)
	err3 := SetTag("refs/tags/v0.1.1", commitOid)

	// then
	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.FileExists(t, ".ugit/refs/tags/v0.1.0")

	h, _ := ioutil.ReadFile(".ugit/refs/tags/v0.1.0")
	assert.Equal(t, "123", string(h))

	assert.Nil(t, err3)
	h, _ = ioutil.ReadFile(".ugit/refs/tags/v0.1.1")
	assert.Equal(t, "123", string(h))
	mock.Teardown()
}

func TestGetTag(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/tags", 0777)
	ioutil.WriteFile(".ugit/refs/tags/v0.1.0", []byte("123"), 0777)

	// when
	oid, _ := GetTag("v0.1.0")
	oid2, _ := GetTag("refs/tags/v0.1.0")

	// then
	assert.Equal(t, "123", oid)
	assert.Equal(t, "123", oid2)
	mock.Teardown()
}

func TestGetBranch(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/heads", 0777)
	ioutil.WriteFile(".ugit/refs/heads/master", []byte("123"), 0777)

	// when
	oid, err := GetBranch("master")
	_, err2 := GetBranch("no_branch")

	// then
	assert.Equal(t, "123", oid)
	assert.Nil(t, err)
	assert.NotNil(t, err2)
	mock.Teardown()
}

func TestSetBranch(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/heads/", 0777)
	ioutil.WriteFile(".ugit/refs/heads/master", []byte(""), 0777)
	commitOid := "123"

	// when
	err1 := SetBranch("master", commitOid)
	err2 := SetBranch("branch_does_not_exist", commitOid)

	// then
	assert.Nil(t, err1)
	assert.FileExists(t, ".ugit/refs/heads/master")
	h, _ := ioutil.ReadFile(".ugit/refs/heads/master")
	assert.Equal(t, "123", string(h))

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
