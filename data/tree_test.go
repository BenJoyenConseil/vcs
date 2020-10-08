package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTmpDir() {
	os.MkdirAll("tmp/other", 0777)
	ioutil.WriteFile("tmp/cats.txt", []byte("Hello"), 0777)
	ioutil.WriteFile("tmp/dogs.txt", []byte("World"), 0777)
	ioutil.WriteFile("tmp/other/shoes.jpg", []byte("qui mange un saucisson"), 0777)
}

func setupUgitDir() {
	os.MkdirAll(".ugit/objects", 0777)
	ioutil.WriteFile(".ugit/objects/429d2f37997444b85323305c5e02c4233a04158e", []byte("blob"+string('\000')+"Hello"), 0777)
	ioutil.WriteFile(".ugit/objects/04921f098f08b8146b16bfdf1173a6cc3013332b", []byte("blob"+string('\000')+"World"), 0777)
	ioutil.WriteFile(".ugit/objects/7a117da734c7e42e7c5a8839715a5a1220a4504f", []byte("blob"+string('\000')+"qui mange un saucisson"), 0777)
	ioutil.WriteFile(".ugit/objects/2099e065ed4f38fc997ca05a706ab6ad31663225", []byte("tree"+string('\000')+"blob 429d2f37997444b85323305c5e02c4233a04158e cats.txt\nblob 04921f098f08b8146b16bfdf1173a6cc3013332b dogs.txt\ntree 2e2df45d8c8bebe3b8945e409f593486ddbc8603 other"), 0777)
	ioutil.WriteFile(".ugit/objects/2e2df45d8c8bebe3b8945e409f593486ddbc8603", []byte("tree"+string('\000')+"blob 7a117da734c7e42e7c5a8839715a5a1220a4504f shoes.jpg"), 0777)
	ioutil.WriteFile(".ugit/objects/323460bfcda38ee6c31f2177e99d7bf1717bf60e", []byte("commit"+string('\000')+"tree 2099e065ed4f38fc997ca05a706ab6ad31663225\nparent \n\nadd something and snapshot it !"), 0777)
	ioutil.WriteFile(".ugit/objects/93584d4997160f16e3ac4390ec4008a2d2ff32d6", []byte("commit"+string('\000')+"tree 2099e065ed4f38fc997ca05a706ab6ad31663225\nparent 323460bfcda38ee6c31f2177e99d7bf1717bf60e\n\nmove you HEAD !"), 0777)
	ioutil.WriteFile(".ugit/objects/cdf776713053cc0710735a61dfbe6492f3ed31b2", []byte("commit"+string('\000')+"tree 2099e065ed4f38fc997ca05a706ab6ad31663225\nparent 93584d4997160f16e3ac4390ec4008a2d2ff32d6\n\nand move again !"), 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("cdf776713053cc0710735a61dfbe6492f3ed31b2"), 0777)
}

func teardown() {
	os.RemoveAll("tmp")
	os.RemoveAll(".ugit")
}

func TestWriteTree(t *testing.T) {
	// given
	setupTmpDir()

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
	assert.Contains(t, string(tmpOtherObj), string(TREE))
	tmpOtherObj, _ = ioutil.ReadFile(".ugit/objects/2099e065ed4f38fc997ca05a706ab6ad31663225")
	assert.Contains(t, string(tmpOtherObj), string(TREE))

	teardown()
}

func TestReadTree(t *testing.T) {

	// given
	oid := "2099e065ed4f38fc997ca05a706ab6ad31663225"
	setupUgitDir()
	setupTmpDir()

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
	teardown()
}

func TestCommit(t *testing.T) {

	//given
	teardown()
	setupTmpDir()

	// when
	oid := Commit("tmp", "add something and snapshot it !")

	// then
	expectedCommitOid := "323460bfcda38ee6c31f2177e99d7bf1717bf60e"
	assert.Equal(t, expectedCommitOid, oid)
	assert.FileExists(t, ".ugit/objects/"+expectedCommitOid)
	assert.FileExists(t, ".ugit/HEAD")
	h, _ := ioutil.ReadFile(".ugit/HEAD")
	assert.Equal(t, expectedCommitOid, string(h))

	c, _ := ioutil.ReadFile(".ugit/objects/" + expectedCommitOid)
	lines := strings.Split(string(bytes.Split(c, []byte{BYTE_SEPARATOR})[1]), "\n")
	assert.Equal(t, "tree 2099e065ed4f38fc997ca05a706ab6ad31663225", lines[0])
	assert.Equal(t, "parent ", lines[1])
	assert.Equal(t, "add something and snapshot it !", lines[3])

	teardown()
}

func TestGetCommit(t *testing.T) {

	// given
	setupUgitDir()
	oid := "93584d4997160f16e3ac4390ec4008a2d2ff32d6"

	// when

	tree, parent, message, err := GetCommit(oid)

	// then

	assert.Equal(t, "2099e065ed4f38fc997ca05a706ab6ad31663225", tree)
	assert.Equal(t, "323460bfcda38ee6c31f2177e99d7bf1717bf60e", parent)
	assert.Equal(t, "move you HEAD !", message)
	assert.Nil(t, err)
}

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
	teardown()
}

func TestGetHead(t *testing.T) {
	// given
	os.MkdirAll(".ugit/", 0777)
	ioutil.WriteFile(".ugit/HEAD", []byte("123"), 0777)

	// when
	oid := GetHead()

	// then
	assert.Equal(t, "123", oid)
	teardown()
}

func TestLog(t *testing.T) {
	// given
	setupUgitDir()

	// when
	commitLog := Log()

	// then
	assert.Contains(t, commitLog.oid, "cdf776713053cc0710735a61dfbe6492f3ed31b2")
	assert.Equal(t, "and move again !", commitLog.message)

	assert.Contains(t, commitLog.parent.oid, "93584d4997160f16e3ac4390ec4008a2d2ff32d6")
	assert.Equal(t, "move you HEAD !", commitLog.parent.message)

	assert.Contains(t, commitLog.parent.parent.oid, "323460bfcda38ee6c31f2177e99d7bf1717bf60e")
	assert.Equal(t, commitLog.parent.parent.message, "add something and snapshot it !")

	teardown()
}
