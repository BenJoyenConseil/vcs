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
	err := SetHead(oid)

	// then
	assert.Nil(t, err)
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
