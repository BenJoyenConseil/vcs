# vcs
![Go](https://github.com/BenJoyenConseil/vcs/workflows/Go/badge.svg)
![Release](https://img.shields.io/github/v/release/BenJoyenConseil/vcs.svg)

Version control system git like

Following the excellent article https://www.leshenko.net/p/ugit/

I have done some modifications
- The `data` module is renammed as the `storage` because it handles every interaction with the filesystem getting and putting objects
- The `base` module is renammed as the `tree` because it contains the hight level structure of a Tree, based on the `storage` module. 
- The code is written using TDD approche. You will find unit tests for each exported functions.

## usage

    go run main.go init | commit | log | ...
  
using binaries
  
    vcs  init | commit | log | ...

## install

    # install using go package manager
    go get github.com/BenJoyenConseil/vcs

    # init to your working directory
    cd my_project
    vcs init

## features

- [x] **commit** -m "add description" : snapshot the current directory and save its version for later
- [x] **hash_object** : add a new object in the datastore folder ./.ugit/objects/
- [x] **log** : print the graph of commit log
- [x] **HEAD** : pointer to the last commit, the current working
- [x] **checkout** : restore files and folders of a specific committed snapshot
- [x] **tag** : mark and identify a commit with a simplified name (instead of its hash)
- [x] **branch** : tag/ref that moves to point on the last commit of the branch
