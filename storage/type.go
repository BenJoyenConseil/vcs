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
	OBJECTS_DIR    string = "objects"
	HEAD_PATH      string = "HEAD"
	REF_DIR        string = "refs"
	TAG_DIR        string = "refs/tags"
	BRANCH_DIR     string = "refs/heads"
	BYTE_SEPARATOR byte   = '\x00'
)
