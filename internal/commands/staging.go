package commands

import (
	"os"
	"slices"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

// Returns the hash of the new blob
func fileToBlob(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	blob := files.Blob(content)
	hash, err := files.CreateObject(&blob)
	if err != nil {
		return "", nil
	}
	return hash, nil
}

func StageFiles(paths []string) error {
	workspaceFiles, err := files.GetWorkspaceFiles()
	if err != nil {
		return err
	}
	indexData, err := files.ReadIndexFile()
	if err != nil {
		return err
	}
	expandedPaths, err := utils.ExpandPaths(paths)
	if err != nil {
		return err
	}
	nonIngoredExpandedPaths := []string{}
	for _, path := range expandedPaths {
		if slices.Contains(workspaceFiles, path) {
			nonIngoredExpandedPaths = append(nonIngoredExpandedPaths, path)
		}
	}
	for _, path := range nonIngoredExpandedPaths {
		hash, err := fileToBlob(path)
		if err != nil {
			return err
		}
		indexData[path] = hash
	}
	err = files.WriteIndexFile(indexData)
	return err
}

func UnstageFiles(paths []string) error {
	indexData, err := files.ReadIndexFile()
	if err != nil {
		return err
	}
	expandedPaths, err := utils.ExpandPaths(paths)
	if err != nil {
		return err
	}
	for _, path := range expandedPaths {
		delete(indexData, path)
	}
	err = files.WriteIndexFile(indexData)
	return err
}
