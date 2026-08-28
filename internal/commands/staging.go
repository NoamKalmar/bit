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
	for _, path := range paths {
		hash, err := fileToBlob(path, hasher)
		if err != nil {
			return err
		}
		stageData[path] = hash
	}
	files.WriteStageFile(stageData)
	return nil
}
