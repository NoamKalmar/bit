package commands

import (
	"path/filepath"
	"strings"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

func CreateCommit(message []string) error {
	commitMessage := strings.Join(message, " ")
	branch, err := files.ReadHeadFile()
	if err != nil {
		return err
	}
	parent, err := files.ReadBranchRef(branch) // last commit is the parent of the new commit
	if err != nil {
		return err
	}
	indexData, err := files.ReadIndexFile()
	if err != nil {
		return err
	}
	treeHash, err := buildTree(".", indexData)
	if err != nil {
		return err
	}
	commit := files.Commit{
		TreeHash:         treeHash,
		ParentCommitHash: parent,
		Message:          commitMessage,
	}
	hash, err := files.CreateObject(&commit)
	if err != nil {
		return err
	}
	files.WriteBranchRef(branch, hash)
	return nil
}

// Returns the hash of the built tree
func buildTree(path string, indexData map[string]string) (string, error) {
	indexedFiles := utils.GetKeys(indexData)
	tree := files.Tree(map[string]string{})
	// Get all of the sub directories to the current directory
	// Build a tree for each one recusivley
	// Add this tree to the current tree
	subdirs, filepaths := utils.ReadDirFromPaths(indexedFiles, path)
	for _, subdir := range subdirs {
		hash, err := buildTree(path+"/"+subdir, indexData)
		if err != nil {
			return "", err
		}
		tree[subdir+"/"] = hash
	}
	// Add blobs to the current tree
	for _, filePath := range filepaths {
		tree[filePath] = indexData[filepath.Join(path, filePath)]
	}
	hash, err := files.CreateObject(&tree)
	if err != nil {
		return "", err
	}
	return hash, nil
}
