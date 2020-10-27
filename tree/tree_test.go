package tree

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"vcs/mock"
	"vcs/storage"

	"github.com/stretchr/testify/assert"
)

func TestWriteTree(t *testing.T) {
	// given
	mock.SetupTmpDir()
	os.MkdirAll(".ugit/objects", 0777)

	// when
	oid, err := WriteTree("./tmp")

	// then
	assert.Nil(t, err)
	assert.Equal(t, oid, "2099e065ed4f38fc997ca05a706ab6ad31663225")
	assert.FileExists(t, ".ugit/objects/2099e065ed4f38fc997ca05a706ab6ad31663225", "/tmp folder does not exists")
	assert.FileExists(t, ".ugit/objects/2e2df45d8c8bebe3b8945e409f593486ddbc8603", "/tmp/other folder does not exists")
	assert.FileExists(t, ".ugit/objects/429d2f37997444b85323305c5e02c4233a04158e", "tmp/cats.txt file does not exists")
	assert.FileExists(t, ".ugit/objects/04921f098f08b8146b16bfdf1173a6cc3013332b", "tmp/dogs.txt file does not exists")
	assert.FileExists(t, ".ugit/objects/7a117da734c7e42e7c5a8839715a5a1220a4504f", "tmp/other/shoes.jpg file does not exists")
	tmpOtherObj, _ := ioutil.ReadFile(".ugit/objects/2e2df45d8c8bebe3b8945e409f593486ddbc8603")
	assert.Contains(t, string(tmpOtherObj), string(storage.TREE))
	tmpOtherObj, _ = ioutil.ReadFile(".ugit/objects/2099e065ed4f38fc997ca05a706ab6ad31663225")
	assert.Contains(t, string(tmpOtherObj), string(storage.TREE))

	mock.Teardown()
}
func TestIsIgnored(t *testing.T) {
	// then
	assert.True(t, IsIgnored("/.ugit"))
	assert.True(t, IsIgnored("/.git"))
	assert.True(t, IsIgnored("tmp/.git"))
	assert.True(t, IsIgnored("/home/lalaland/.local/.git/HEAD"))
	assert.False(t, IsIgnored("/home/lalaland/.local/yo"))
}

func TestReadTree(t *testing.T) {

	// given
	oid := "2099e065ed4f38fc997ca05a706ab6ad31663225"
	mock.SetupUgitDir()

	// // when
	err := ReadTree(oid, "tmp")

	// // then
	assert.Nil(t, err)
	assert.DirExists(t, "tmp")
	assert.DirExists(t, "tmp/other")
	assert.FileExists(t, "tmp/cats.txt")
	assert.FileExists(t, "tmp/dogs.txt")
	assert.FileExists(t, "tmp/other/shoes.jpg")

	d, _ := ioutil.ReadFile("tmp/other/shoes.jpg")
	assert.Equal(t, "qui mange un saucisson", string(d))
	d, _ = ioutil.ReadFile("tmp/cats.txt")
	assert.Equal(t, "Hello", string(d))
	d, _ = ioutil.ReadFile("tmp/dogs.txt")
	assert.Equal(t, "World", string(d))
	mock.Teardown()
}

func TestReadTree_ShouldErrorIfObjectNotATree(t *testing.T) {
	// given
	mock.SetupUgitDir()
	oidCommitMoveAgain := "cdf776713053cc0710735a61dfbe6492f3ed31b2"

	// when
	err := ReadTree(oidCommitMoveAgain, "tmp")

	// then
	assert.NotNil(t, err)
	mock.Teardown()
}

func TestReadTree_ShouldNotRemoveIgnoredPath(t *testing.T) {
	// given
	mock.SetupUgitDir()
	mock.SetupTmpDir()
	ioutil.WriteFile("tmp/.gitignore", []byte(""), 0777)
	os.MkdirAll("tmp/.git", 0777)
	os.MkdirAll("tmp/.ugit", 0777)
	oidCommitMoveAgain := "cdf776713053cc0710735a61dfbe6492f3ed31b2"

	// when
	ReadTree(oidCommitMoveAgain, "tmp")

	// then
	assert.FileExists(t, "tmp/.gitignore")
	assert.DirExists(t, "tmp/.ugit")
	assert.DirExists(t, "tmp/.git")
	mock.Teardown()
}

func TestCommit(t *testing.T) {

	//given
	mock.SetupTmpDir()
	os.MkdirAll(".ugit/refs/heads", 0777)
	os.MkdirAll(".ugit/objects", 0777)
	ioutil.WriteFile(".ugit/refs/heads/master", nil, 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("refs/heads/master"), 0777)

	// when
	oid := Commit("tmp", "add something and snapshot it !")

	// then
	expectedCommitOid := "323460bfcda38ee6c31f2177e99d7bf1717bf60e"
	assert.Equal(t, expectedCommitOid, oid)
	assert.FileExists(t, ".ugit/objects/"+expectedCommitOid)

	h, _ := ioutil.ReadFile(".ugit/refs/heads/master")
	assert.Equal(t, expectedCommitOid, string(h))

	c, _ := ioutil.ReadFile(".ugit/objects/" + expectedCommitOid)
	lines := strings.Split(string(bytes.Split(c, []byte{storage.BYTE_SEPARATOR})[1]), "\n")
	assert.Equal(t, "tree 2099e065ed4f38fc997ca05a706ab6ad31663225", lines[0])
	assert.Equal(t, "parent ", lines[1])
	assert.Equal(t, "add something and snapshot it !", lines[3])

	mock.Teardown()
}

func TestGetCommit(t *testing.T) {

	// given
	mock.SetupUgitDir()
	oid := "93584d4997160f16e3ac4390ec4008a2d2ff32d6"

	// when

	tree, parent, message, err := GetCommit(oid)

	// then

	assert.Equal(t, "2099e065ed4f38fc997ca05a706ab6ad31663225", tree)
	assert.Equal(t, "323460bfcda38ee6c31f2177e99d7bf1717bf60e", parent)
	assert.Equal(t, "move you HEAD !", message)
	assert.Nil(t, err)
	mock.Teardown()
}

func TestLog(t *testing.T) {
	// given
	mock.SetupUgitDir()

	// when
	commitLog := Log()

	// then
	assert.Contains(t, commitLog.oid, "cdf776713053cc0710735a61dfbe6492f3ed31b2")
	assert.Equal(t, "and move again !", commitLog.message)

	assert.Contains(t, commitLog.parent.oid, "93584d4997160f16e3ac4390ec4008a2d2ff32d6")
	assert.Equal(t, "move you HEAD !", commitLog.parent.message)

	assert.Contains(t, commitLog.parent.parent.oid, "323460bfcda38ee6c31f2177e99d7bf1717bf60e")
	assert.Equal(t, commitLog.parent.parent.message, "add something and snapshot it !")

	mock.Teardown()
}

func TestCheckout(t *testing.T) {
	// given
	mock.SetupTmpDir()
	mock.SetupUgitDir()
	mock.RemoveDogsAndCommit()
	oid := "f37333b2d9ffbbf083b6c364a02cc555fa56ffef"

	// when
	err := Checkout(oid, "tmp")

	// then
	assert.Nil(t, err)
	assert.DirExists(t, "tmp")
	assert.FileExists(t, "tmp/cats.txt")
	assert.DirExists(t, "tmp/other")
	assert.FileExists(t, "tmp/other/shoes.jpg")
	assert.NoFileExists(t, "tmp/dogs.txt")
	h, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "f37333b2d9ffbbf083b6c364a02cc555fa56ffef", string(h))
	mock.Teardown()
}
func TestCheckout_Tag(t *testing.T) {
	// given
	mock.SetupTmpDir()
	mock.SetupUgitDir()
	mock.RemoveDogsAndCommit()
	os.MkdirAll(".ugit/refs/tags", 0777)
	rmDogCommitOid := []byte("f37333b2d9ffbbf083b6c364a02cc555fa56ffef")
	ioutil.WriteFile(".ugit/refs/tags/v0.1.0", rmDogCommitOid, 0777)

	// when
	err := Checkout("v0.1.0", "tmp")

	// then
	assert.Nil(t, err)
	assert.DirExists(t, "tmp")
	assert.FileExists(t, "tmp/cats.txt")
	assert.DirExists(t, "tmp/other")
	assert.FileExists(t, "tmp/other/shoes.jpg")
	assert.NoFileExists(t, "tmp/dogs.txt")

	h, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, "f37333b2d9ffbbf083b6c364a02cc555fa56ffef", string(h))

	mock.Teardown()
}

func TestGetOid(t *testing.T) {
	mock.SetupUgitDir()

	// when tag
	oid, fullref := GetOid("v0.1.0")
	oid2, fullref2 := GetOid("refs/tags/v0.1.0")
	// then
	assert.Equal(t, "93584d4997160f16e3ac4390ec4008a2d2ff32d6", oid)
	assert.Equal(t, "refs/tags/v0.1.0", fullref)
	assert.Equal(t, "93584d4997160f16e3ac4390ec4008a2d2ff32d6", oid2)
	assert.Equal(t, "refs/tags/v0.1.0", fullref2)

	// when commit
	oid, fullref = GetOid("323460bfcda38ee6c31f2177e99d7bf1717bf60e")
	// then
	assert.Equal(t, "323460bfcda38ee6c31f2177e99d7bf1717bf60e", oid)
	assert.Equal(t, "323460bfcda38ee6c31f2177e99d7bf1717bf60e", fullref)

	// when branch
	oid, fullref = GetOid("master")
	oid2, fullref2 = GetOid("refs/heads/master")
	// then
	assert.Equal(t, "cdf776713053cc0710735a61dfbe6492f3ed31b2", oid)
	assert.Equal(t, "refs/heads/master", fullref)
	assert.Equal(t, "cdf776713053cc0710735a61dfbe6492f3ed31b2", oid2)
	assert.Equal(t, "refs/heads/master", fullref2)

	// when custom ref
	oid, fullref = GetOid("HEAD")
	// then
	assert.Equal(t, "refs/heads/master", oid)
	assert.Equal(t, "HEAD", fullref)

	mock.Teardown()
}

func TestCreateTag(t *testing.T) {
	// given
	os.MkdirAll(".ugit/refs/tags", 0777)
	commitOid := "123"

	// when
	noErr := CreateTag("v0.1.0", commitOid)
	fail := CreateTag("v0.1.0", commitOid)

	// then
	assert.Nil(t, noErr)
	assert.NotNil(t, fail)
	assert.FileExists(t, ".ugit/refs/tags/v0.1.0")

	h, _ := ioutil.ReadFile(".ugit/refs/tags/v0.1.0")
	assert.Equal(t, "123", string(h))
	mock.Teardown()
}
