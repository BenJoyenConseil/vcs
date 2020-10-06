package data

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteTree(t *testing.T) {
	// given
	os.MkdirAll("tmp/other", 0777)
	ioutil.WriteFile("tmp/cats.txt", []byte("Hello"), 0777)
	ioutil.WriteFile("tmp/dogs.txt", []byte("World"), 0777)
	ioutil.WriteFile("tmp/other/shoes.jpg", []byte("qui mange un saucisson"), 0777)

	// when
	oid, err := WriteTree("./tmp")

	// then
	resultFiles, _ := ioutil.ReadDir(".ugit/objects/")
	expectedFiles := map[string]string{
		"429d2f37997444b85323305c5e02c4233a04158e": "tmp/cats.txt",
		"04921f098f08b8146b16bfdf1173a6cc3013332b": "tmp/dogs.txt",
		"7a117da734c7e42e7c5a8839715a5a1220a4504f": "tmp/other/shoes.jpg",
		"2099e065ed4f38fc997ca05a706ab6ad31663225": "tmp",
		"2e2df45d8c8bebe3b8945e409f593486ddbc8603": "tmp/other/",
	}
	if err != nil {
		t.Error(err)
	}
	if oid != "2099e065ed4f38fc997ca05a706ab6ad31663225" {
		t.Error("oid not match ", expectedFiles["2099e065ed4f38fc997ca05a706ab6ad31663225"], oid)
	}
	for _, f := range resultFiles {
		if val, ok := expectedFiles[f.Name()]; ok {
			t.Log("Found file : ", val, " for ", f.Name())
		} else {
			t.Error(f.Name(), " is not expected in")
		}
	}

	// teardown
	os.RemoveAll("tmp")
	os.RemoveAll(".ugit")
}
