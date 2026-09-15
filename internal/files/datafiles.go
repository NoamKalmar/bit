package files

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/noamkalmar/bit/internal/utils"
)

func ReadIndexFile() (map[string]string, error) {
	return utils.ReadJsonStrStr(INDEX_PATH)
}

func WriteIndexFile(stagingData map[string]string) error {
	return utils.WriteJsonStrStr(INDEX_PATH, stagingData)
}

func ReadHeadFile() (string, error) {
	content, err := os.ReadFile(HEAD_PATH)
	return string(content), err
}

func WriteHeadFile(branchName string) error {
	return os.WriteFile(HEAD_PATH, []byte(branchName), 0644)
}

func IgnoreFileExists() bool {
	return utils.IsFile(IGNORE_PATH)
}

func ReadIgnoreFile() ([]string, error) {
	content, err := os.ReadFile(IGNORE_PATH)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(string(content), "\n")
	for i, path := range paths {
		paths[i] = filepath.Clean(strings.TrimSpace(path))
	}
	return paths, nil
}

func GetAllNonIgnoredFiles() ([]string, error) {
	ignorePaths := []string{}
	if IgnoreFileExists() {
		var err error
		ignorePaths, err = ReadIgnoreFile()
		if err != nil {
			return nil, err
		}
	}
	ignorePaths = append(ignorePaths, BIT_PATH) // the .bit folder should be ignored
	ignoreFiles, ignoreDirs := utils.ClassifyPaths(ignorePaths)

	files, err := utils.GetAllFilePaths()
	if err != nil {
		return nil, err
	}

	// removing all files that appear in the ignore file
	files = slices.DeleteFunc(files, func(path string) bool {
		return slices.Contains(ignoreFiles, path)
	})

	// removing all files that are in the tree of an ignored directory
	files = slices.DeleteFunc(files, func(path string) bool {
		for _, dir := range ignoreDirs {
			if strings.HasPrefix(path, dir+"\\") {
				return true
			}
		}
		return false
	})
	return files, nil
}
