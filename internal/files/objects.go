package files

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/noamkalmar/bit/internal/utils"
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

type Object interface {
	Type() ObjectType

	Serialize() ([]byte, error)
}

type Blob []byte

type Tree map[string]string

// A commit object has a defined structure, therefore we can define how to parse it as json
type Commit struct {
	TreeHash         string `json:"tree"`
	ParentCommitHash string `json:"parent"`
	Message          string `json:"message"`
}

func (b *Blob) Type() ObjectType {
	return BlobType
}

func (b *Blob) Serialize() ([]byte, error) {
	return []byte(*b), nil
}

func (t *Tree) Type() ObjectType {
	return TreeType
}

func (t *Tree) Serialize() ([]byte, error) {
	return json.MarshalIndent(*t, "", "    ")
}

func (c *Commit) Type() ObjectType {
	return CommitType
}

func (c *Commit) Serialize() ([]byte, error) {
	return json.MarshalIndent(*c, "", "    ")
}

// Returns the hash of the object
func CreateObject(object Object) (string, error) {
	folderName := objectTypeToFolderName[object.Type()]
	content, err := object.Serialize()
	if err != nil {
		return "", err
	}
	hash := utils.GetSHA1(content)
	file, err := os.Create(filepath.Join(OBJECTS_PATH, folderName, hash))
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(content)
	return hash, nil
}

func ReadCommit(hash string) (Commit, error) {
	content, err := os.ReadFile(filepath.Join(COMMITS_PATH, hash))
	if err != nil {
		return Commit{}, nil
	}
	var commit Commit
	err = json.Unmarshal(content, &commit)
	if err != nil {
		return Commit{}, nil
	}
	return commit, nil
}
