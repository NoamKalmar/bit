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
	untracked, _, _, err := classfilyFilesByIndex()
	if err != nil {
		return err
	}
	fmt.Println("\nUntracked files: ")
	for _, file := range untracked {
		fmt.Println(" - " + file)
	}
	return nil
}

// Returns untracked, modified files, deleted files
func classfilyFilesByIndex() ([]string, []string, []string, error) {
	allFiles, err := files.GetAllNonIgnoredFiles()
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
