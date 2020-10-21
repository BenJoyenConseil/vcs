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
	oid := "123"
	os.MkdirAll(".ugit/", 0777)

	// when
	SetHead(oid)

	// then
	assert.FileExists(t, ".ugit/HEAD")
	h, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "123", string(h))
	mock.Teardown()
}

func TestGetHead(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("123"), 0777)

	// when
	oid := GetHead()

	// then
	assert.Equal(t, "123", oid)
	mock.Teardown()
}

func TestSetTag(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)
	commitOid := "123"

	// when
	err1 := SetTag("v0.1.0", commitOid)
	err2 := SetTag("v0.1.0", commitOid)

	// then
	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.FileExists(t, ".ugit/refs/tags/v0.1.0")
	h, _ := ioutil.ReadFile(".ugit/refs/tags/v0.1.0")
	assert.Equal(t, "123", string(h))
	mock.Teardown()
}

func TestGetTag(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/tags", 0777)
	ioutil.WriteFile(".ugit/refs/tags/v0.1.0", []byte("123"), 0777)

	// when
	oid := GetTag("/v0.1.0")

	// then
	assert.Equal(t, "123", oid)
	mock.Teardown()
}
