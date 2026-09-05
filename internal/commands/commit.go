package commands

import (
	"crypto/sha1"
	"encoding/json"
	"hash"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

func CreateCommit(message []string) error {
	commitMessage := strings.Join(message, "")
	branch, err := files.ReadHeadFile()
	if err != nil {
		return err
	}
	parent, err := files.ReadBranchRef(branch) // last commit is the parent of the new commit
	if err != nil {
		return err
	}
	stageData, err := files.ReadStageFile()
	if err != nil {
		return err
	}
	hasher := sha1.New()
	buildTree(".", stageData, hasher)
	commit := files.Commit{
		TreeHash:         "",
		ParentCommitHash: parent,
		Message:          commitMessage,
	}
	content, err := json.MarshalIndent(commit, "", "    ")
	if err != nil {
		return err
	}
	hash := utils.GetHash(content, hasher)
	files.CreateCommit(hash, content) // Create a new commit object
	files.WriteBranchRef(branch, hash)
	return nil
}

// Returns the hash of the built tree
func buildTree(path string, stageData map[string]string, hasher hash.Hash) (string, error) {
	toStage := slices.Collect(maps.Keys(stageData))
	tree := map[string]string{}
	// Get all of the sub directories to the current directory
	// Build a tree for each one recusivley
	// Add this tree to the current tree
	subdirs, filepaths := utils.ReadDirFromPaths(toStage, path)
	for _, subdir := range subdirs {
		hash, err := buildTree(path+"/"+subdir, stageData, hasher)
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
	hash := utils.GetHash(content, hasher)
	files.CreateTree(hash, content)
	return hash, nil
}
