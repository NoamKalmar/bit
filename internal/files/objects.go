package files

import (
	"os"
	"path/filepath"
)

type ObjectType int

const (
	BlobType = iota
	TreeType
	CommitType
)

var objectTypeToFolderName = map[ObjectType]string{
	BlobType:   "blobs",
	TreeType:   "trees",
	CommitType: "commits",
}

// A commit object has a defined structure, therefore we can define how to parse it as json
type Commit struct {
	TreeHash         string `json:"tree"`
	ParentCommitHash string `json:"parent"`
	Message          string `json:"message"`
}

func createObject(objectType ObjectType, hash string, content []byte) error {
	folderName := objectTypeToFolderName[objectType]
	file, err := os.Create(filepath.Join(OBJECTS_PATH, folderName, hash))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(content)
	return nil
}

func CreateBlob(hash string, content []byte) error {
	return createObject(BlobType, hash, content)
}

func CreateCommit(hash string, content []byte) error {
	return createObject(CommitType, hash, content)
}

func CreateTree(hash string, content []byte) error {
	return createObject(TreeType, hash, content)
}
