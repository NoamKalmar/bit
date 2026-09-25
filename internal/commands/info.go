package commands

import (
	"fmt"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

// Options when checking last commit:
// modified (hash != old hash)
// new file (appears in workspace but not in commit)
// deleted (appears in commit but not in workspace)

// Then for each one, check if the stage file is updated for it.
// If it is, add it to the changes to be commited
// If it's not, add it to the changes not staged for commit

func PrintStatus() error {
	branch, err := files.ReadHeadFile()
	if err != nil {
		return err
	}
	fmt.Println("Current branch: " + branch)
	modified, new, deleted, err := getDiffFromLastCommit()
	staged, nonStaged, err := getStagedAndNotStagedChanges()
	if err != nil {
		return err
	}

	fmt.Println("\nChanges to be commited: ")
	for file := range staged {
		if modified.Contains(file) {
			fmt.Println(" - modified: " + file)
		} else if new.Contains(file) {
			fmt.Println(" - new: " + file)
		} else if deleted.Contains(file) {
			fmt.Println(" - deleted: " + file)
		}
	}

	fmt.Println("\nChanges not staged for commit: ")
	for file := range nonStaged {
		if modified.Contains(file) {
			fmt.Println(" - modified: " + file)
		} else if deleted.Contains(file) {
			fmt.Println(" - deleted: " + file)
		}
	}

	fmt.Println("\nUntracked files: ")
	for file := range nonStaged {
		if new.Contains(file) {
			fmt.Println(" - " + file)
		}
	}

	return nil
}

// Returns modified, new, deleted
func getDiffFromLastCommit() (utils.PathSet, utils.PathSet, utils.PathSet, error) {
	commit, err := files.GetLastCommit()
	if err != nil {
		return nil, nil, nil, err
	}
	tree, err := files.ReadTree(commit.TreeHash)
	if err != nil {
		return nil, nil, nil, err
	}
	commitFiles, err := tree.Walk()
	if err != nil {
		return nil, nil, nil, err
	}
	commitFilePaths := utils.SliceToPathSet(utils.GetKeys(commitFiles))

	workspaceFilePaths, err := files.GetWorkspaceFiles()
	if err != nil {
		return nil, nil, nil, err
	}
	// All files that appear in the workspace but not in the commit are new files
	// All files that appear in the commit but not in the workspace are deleted files
	new := workspaceFilePaths.Difference(commitFilePaths)
	deleted := commitFilePaths.Difference(workspaceFilePaths)

	// Modified files have a different hash than the one they had in the last commit
	modified := make(utils.PathSet)
	for path := range workspaceFilePaths {
		if !commitFilePaths.Contains(path) {
			continue
		}
		oldHash := commitFiles[path]
		hash, err := utils.ReadGetSHA1(path)
		if err != nil {
			return nil, nil, nil, err
		}
		if hash != oldHash {
			modified.Add(path)
		}
	}

	return modified, new, deleted, nil
}

// Returns staged, not staged
func getStagedAndNotStagedChanges() (utils.PathSet, utils.PathSet, error) {
	notStaged := make(utils.PathSet)
	indexData, err := files.ReadIndexFile()
	workspaceFiles, err := files.GetWorkspaceFiles()
	if err != nil {
		return nil, nil, err
	}
	indexFiles := utils.SliceToPathSet(utils.GetKeys(indexData))

	new := workspaceFiles.Difference(indexFiles)
	deleted := indexFiles.Difference(workspaceFiles)

	// either new or deleted files (according to the index) are not staged for commit
	notStaged.AddSet(new)
	notStaged.AddSet(deleted)
	for path := range workspaceFiles {
		if !indexFiles.Contains(path) {
			continue
		}
		oldHash := indexData[path]
		hash, err := utils.ReadGetSHA1(path)
		if err != nil {
			return nil, nil, err
		}
		if oldHash != hash {
			notStaged.Add(path)
		}
	}
	staged := workspaceFiles.Difference(notStaged)
	return staged, notStaged, nil
}
