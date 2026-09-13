package commands

import (
	"encoding/json"
	"maps"
	"path/filepath"
	"slices"
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
	stageData, err := files.ReadIndexFile()
	if err != nil {
		return err
	}
	treeHash, err := buildTree(".", stageData)
	if err != nil {
		return err
	}
	commit := files.Commit{
		TreeHash:         treeHash,
		ParentCommitHash: parent,
		Message:          commitMessage,
	}
	content, err := json.MarshalIndent(commit, "", "    ")
	if err != nil {
		return err
	}
	hash := utils.GetSHA1(content)
	files.CreateCommit(hash, content) // Create a new commit object
	files.WriteBranchRef(branch, hash)
	return nil
}

// Returns the hash of the built tree
func buildTree(path string, stageData map[string]string) (string, error) {
	toStage := slices.Collect(maps.Keys(stageData))
	tree := map[string]string{}
	// Get all of the sub directories to the current directory
	// Build a tree for each one recusivley
	// Add this tree to the current tree
	subdirs, filepaths := utils.ReadDirFromPaths(toStage, path)
	for _, subdir := range subdirs {
		hash, err := buildTree(path+"/"+subdir, stageData)
		if err != nil {
			return "", err
		}
		tree[subdir+"/"] = hash
	}
	// Add blobs to the current tree
	for _, filePath := range filepaths {
		tree[filePath] = stageData[filepath.Join(path, filePath)]
	}
	content, err := json.MarshalIndent(tree, "", "    ")
	if err != nil {
		return "", err
	}
	hash := utils.GetSHA1(content)
	files.CreateTree(hash, content)
	return hash, nil
}
