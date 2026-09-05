package files

import (
	"os"
)

// A commit object has a defined structure, therefore we can define how to parse it as json
type Commit struct {
	TreeHash         string `json:"tree"`
	ParentCommitHash string `json:"parent"`
	Message          string `json:"message"`
}

func createObject(objectType string, hash string, content []byte) error {
	file, err := os.Create(".bit/objects/" + objectType + "/" + hash)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(content)
	return nil
}

func CreateBlob(hash string, content []byte) error {
	return createObject("blobs", hash, content)
}

func CreateCommit(hash string, content []byte) error {
	return createObject("commits", hash, content)
}

func CreateTree(hash string, content []byte) error {
	return createObject("trees", hash, content)
}
