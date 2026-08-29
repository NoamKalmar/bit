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
	hash := utils.GetHash(content, hasher)
	files.CreateBlob(hash, content)
	return hash, nil
}

func StageFiles(paths []string) error {
	stageData, err := files.ReadStageFile()
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
		stageData[path] = hash
	}
	err = files.WriteStageFile(stageData)
	return err
}

func UnstageFiles(paths []string) error {
	stageData, err := files.ReadStageFile()
	if err != nil {
		return err
	}
	expandedPaths, err := utils.ExpandPaths(paths)
	if err != nil {
		return err
	}
	for _, path := range expandedPaths {
		delete(stageData, path)
	}
	err = files.WriteStageFile(stageData)
	return err
}
