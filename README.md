# vcs
![Go](https://github.com/BenJoyenConseil/vcs/workflows/Go/badge.svg)

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
