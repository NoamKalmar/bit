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

func GetLastCommit() (Commit, error) {
	hash, err := GetLastCommitHash()
	if err != nil {
		return Commit{}, nil
	}
	return ReadCommit(hash)
}

func ReadTree(hash string) (Tree, error) {
	content, err := os.ReadFile(filepath.Join(TREES_PATH, hash))
	if err != nil {
		return nil, nil
	}
	var tree map[string]string
	err = json.Unmarshal(content, &tree)
	if err != nil {
		return nil, nil
	}
	return tree, nil
}

// Returns map of file paths to file hashes, map of sub tree path to tree object
func (t *Tree) Classify() (map[string]string, map[string]Tree, error) {
	files := map[string]string{}
	trees := map[string]Tree{}
	for path, hash := range *t {
		if path[len(path)-1] == '/' {
			tree, err := ReadTree(hash)
			if err != nil {
				return nil, nil, err
			}
			trees[path[:len(path)-1]] = tree
		} else {
			files[path] = hash
		}
	}
	return files, trees, nil
}

// Returns all the file paths that appear in some tree and their hashes
func (t *Tree) Walk() (map[string]string, error) {
	files := map[string]string{}
	files, err := t.walk(files, "")
	return files, err
}

func (t *Tree) walk(files map[string]string, parentTreePath string) (map[string]string, error) {
	treeFiles, trees, err := t.Classify()
	if err != nil {
		return nil, err
	}
	for filePath, hash := range treeFiles {
		files[filepath.Join(parentTreePath, filePath)] = hash
	}
	for treePath, tree := range trees {
		files, err = tree.walk(files, filepath.Join(parentTreePath, treePath))
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
