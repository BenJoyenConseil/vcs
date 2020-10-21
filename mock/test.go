package mock

import (
	"io/ioutil"
	"os"
)

func SetupTmpDir() {
	os.MkdirAll("tmp/other", 0777)
	ioutil.WriteFile("tmp/cats.txt", []byte("Hello"), 0777)
	ioutil.WriteFile("tmp/dogs.txt", []byte("World"), 0777)
	ioutil.WriteFile("tmp/other/shoes.jpg", []byte("qui mange un saucisson"), 0777)
}

func SetupUgitDir() {
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

func RemoveDogsAndCommit() {
	ioutil.WriteFile(".ugit/objects/751eb07c9a747033c359510dd71c8dd045a9cfc1", []byte("tree"+string('\000')+"blob 429d2f37997444b85323305c5e02c4233a04158e cats.txt\ntree 2e2df45d8c8bebe3b8945e409f593486ddbc8603 other"), 0777)
	ioutil.WriteFile(".ugit/objects/f37333b2d9ffbbf083b6c364a02cc555fa56ffef", []byte("commit"+string('\000')+"tree 751eb07c9a747033c359510dd71c8dd045a9cfc1\nparent cdf776713053cc0710735a61dfbe6492f3ed31b2\n\nremove dogs.txt"), 0777)
}

func Teardown() {
	os.RemoveAll("tmp")
	os.RemoveAll(".ugit")
}
