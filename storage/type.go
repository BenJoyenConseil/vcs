package storage

type ObjectType string

const (
	BLOB   ObjectType = "blob"
	TREE   ObjectType = "tree"
	COMMIT ObjectType = "commit"
	PARENT ObjectType = "parent"
)

const (
	UGIT_DIR       string = ".ugit"
	OBJECTS_DIR    string = ".ugit/objects"
	HEAD_PATH      string = ".ugit/HEAD"
	REF_DIR        string = ".ugit/refs"
	TAG_DIR        string = ".ugit/refs/tags"
	BRANCH_DIR     string = ".ugit/refs/heads"
	BYTE_SEPARATOR byte   = '\x00'
)
