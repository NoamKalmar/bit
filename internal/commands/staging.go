package commands

import (
	"crypto/sha1"
	"hash"
	"os"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

// Returns the hash of the new blob
func fileToBlob(path string, hasher hash.Hash) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := utils.GetSHA1(content)
	files.CreateBlob(hash, content)
	return hash, nil
}

func StageFiles(paths []string) error {
	indexData, err := files.ReadIndexFile()
	if err != nil {
		return err
	}
	hasher := sha1.New()
	expandedPaths, err := utils.ExpandPaths(paths)
	if err != nil {
		return err
	}
	for _, path := range expandedPaths {
		hash, err := fileToBlob(path, hasher)
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
