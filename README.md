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

    vcs [<flags>] <command> [<args> ...]

    Snapshot your working directory

    Flags:
      --help  Show context-sensitive help (also try --help-long and --help-man).

    Commands:
      help [<command>...]
        Show help.

      commit --message=MESSAGE [<flags>]
        snapshot the current directory with an explicite message desciption

      init [<dir>]
        Setup the directory you want to be managed

      log [<ref>]
        Print the commit log history

      hash_object [<flags>] <content>
        Save an object in vcs and get its hash

      checkout [<flags>] <oid>
        Restore files and folders from a committed snapshot

      tag <name>
        Give a name to the current commit

      branch
        Print all created branches



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
- [ ] **graph** : print branches history and tags with HEAD position
- [ ] **diff** : see the differences of content between 2 working trees
- [ ] **status** : show the current branch, the staging files
- [ ] **merge** : merge content of two branches and create a merge commit with 2 parents
- [ ] **fast-forward** : during a merge, see if HEAD is a ancestor of the head commit of the branch
- [ ] **index** : commit only tracked files
- [ ] **stash** : push pending changes to a stack and remove them of the index
- [ ] **cherry-pick** : pick a commit from a branch and add it to another branch
- [ ] **rebase** : play the branch history to another branch history, and get it linear
