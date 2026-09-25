package files

import (
	"os"
	"strings"

	"github.com/noamkalmar/bit/internal/utils"
)

func ReadIndexFile() (map[string]string, error) {
	return utils.ReadJsonStrStr(INDEX_PATH)
}

func WriteIndexFile(indexData map[string]string) error {
	return utils.WriteJsonStrStr(INDEX_PATH, indexData)
}

func ReadHeadFile() (string, error) {
	content, err := os.ReadFile(HEAD_PATH)
	return string(content), err
}

func WriteHeadFile(branchName string) error {
	return os.WriteFile(HEAD_PATH, []byte(branchName), 0644)
}

func GetLastCommitHash() (string, error) {
	branch, err := ReadHeadFile()
	if err != nil {
		return "", err
	}
	commit, err := ReadBranchRef(branch)
	return commit, err
}

func IgnoreFileExists() bool {
	return utils.IsFile(IGNORE_PATH)
}

func ReadIgnoreFile() (utils.PathSet, error) {
	content, err := os.ReadFile(IGNORE_PATH)
	if err != nil {
		return nil, err
	}
	pathSlice := strings.Split(string(content), "\n")
	paths := make(utils.PathSet)
	for _, path := range pathSlice {
		paths.Add(strings.TrimSpace(path))
	}
	return paths, nil
}

func GetWorkspaceFiles() (utils.PathSet, error) {
	ignorePaths := make(utils.PathSet)
	if IgnoreFileExists() {
		var err error
		ignorePaths, err = ReadIgnoreFile()
		if err != nil {
			return nil, err
		}
	}
	ignorePaths.Add(BIT_PATH) // the .bit folder should be ignored
	ignoreFiles, ignoreDirs := utils.ClassifyPaths(ignorePaths)

	files, err := utils.GetAllFilePaths()
	if err != nil {
		return nil, err
	}

	// removing all files that appear in the ignore file
	for path := range files {
		if ignoreFiles.Contains(path) {
			files.Remove(path)
		}
	}

	// removing all files that are in the tree of an ignored directory
	for path := range files {
		for dir := range ignoreDirs {
			if strings.HasPrefix(path, dir+"\\") || strings.HasPrefix(path, dir+"/") {
				files.Remove(path)
			}
		}
	}
	return files, nil
}
