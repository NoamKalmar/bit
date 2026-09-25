package commands

import (
	"fmt"
	"os"
	"slices"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

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
	// Modified + non staged = not commit: modified
	// Modified + staged = to commit: modified
	// New + non staged = untracked
	// New + staged = to commit: new file
	// Deleted + non staged = not commit: deleted
	// Deleted + staged = to commit: deleted

	fmt.Println("\nChanges to be commited: ")
	for _, file := range staged {
		if slices.Contains(modified, file) {
			fmt.Println(" - modified: " + file)
		} else if slices.Contains(new, file) {
			fmt.Println(" - new: " + file)
		} else if slices.Contains(deleted, file) {
			fmt.Println(" - deleted: " + file)
		}
	}

	fmt.Println("\nChanges not staged for commit: ")
	for _, file := range nonStaged {
		if slices.Contains(modified, file) {
			fmt.Println(" - modified: " + file)
		} else if slices.Contains(deleted, file) {
			fmt.Println(" - deleted: " + file)
		}
	}

	fmt.Println("\nUntracked files: ")
	for _, file := range nonStaged {
		if slices.Contains(new, file) {
			fmt.Println(" - " + file)
		}
	}

	return nil
}

// Options when checking last commit:
// modified (hash != old hash)
// new file (appears in workspace but not in commit)
// deleted (appears in commit but not in workspace)

// Then for each one, check if the stage file is updated for it.
// If it is, add it to the changes to be commited
// If it's not, add it to the changes not staged for commit

// Returns modified, new, deleted
func getDiffFromLastCommit() ([]string, []string, []string, error) {

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
	commitFilePaths := utils.GetKeys(commitFiles)

	workspaceFilePaths, err := files.GetWorkspaceFiles()
	if err != nil {
		return nil, nil, nil, err
	}
	// All files that appear in the workspace but not in the commit are new files
	// All files that appear in the commit but not in the workspace are deleted files
	new, deleted := utils.SlicesDifferences(workspaceFilePaths, commitFilePaths)

	// Modified files have a different hash than the one they had in the last commit
	modified := []string{}
	for _, path := range workspaceFilePaths {
		if !slices.Contains(commitFilePaths, path) {
			continue
		}
		oldHash := commitFiles[path]
		hash, err := utils.ReadGetSHA1(path)
		if err != nil {
			return nil, nil, nil, err
		}
		if hash != oldHash {
			modified = append(modified, path)
		}
	}

	return modified, new, deleted, nil
}

// Returns staged, not staged
func getStagedAndNotStagedChanges() ([]string, []string, error) {
	staged := []string{}
	notStaged := []string{}
	indexData, err := files.ReadIndexFile()
	workspaceFiles, err := files.GetWorkspaceFiles()
	if err != nil {
		return nil, nil, err
	}
	indexFiles := utils.GetKeys(indexData)

	new, deleted := utils.SlicesDifferences(workspaceFiles, indexFiles)
	// either new or deleted files (according to the index) are not staged for commit
	notStaged = append(new, deleted...)
	for _, path := range workspaceFiles {
		if !slices.Contains(indexFiles, path) {
			continue
		}
		oldHash := indexData[path]
		hash, err := utils.ReadGetSHA1(path)
		if err != nil {
			return nil, nil, err
		}
		if oldHash != hash {
			notStaged = append(notStaged, path)
		}
	}
	for _, file := range workspaceFiles {
		if !slices.Contains(notStaged, file) {
			staged = append(staged, file)
		}
	}
	return staged, notStaged, nil
}

func classfilyFilesByIndex() ([]string, []string, []string, error) {
	allFiles, err := files.GetWorkspaceFiles()
	if err != nil {
		return nil, nil, nil, err
	}
	indexData, err := files.ReadIndexFile()
	indexedFiles := utils.GetKeys(indexData)
	if err != nil {
		return nil, nil, nil, err
	}

	untracked := []string{}
	modified := []string{}
	deleted := []string{}

	for _, file := range allFiles {
		if !slices.Contains(indexedFiles, file) {
			untracked = append(untracked, file)
		} else {
			content, err := os.ReadFile(file)
			if err != nil {
				return nil, nil, nil, err
			}
			currentHash := utils.GetSHA1(content)
			oldHash := indexData[file]
			if currentHash != oldHash {
				modified = append(modified, file)
			}
		}
	}
	return untracked, modified, deleted, nil
}
