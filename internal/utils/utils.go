package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func GetSHA1(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func ReadGetSHA1(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	return GetSHA1(content), nil
}

func CreateDirs(paths []string) error {
	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	return false
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// Returns files, dirs
func ClassifyPaths(paths []string) ([]string, []string) {
	files := []string{}
	dirs := []string{}
	for _, path := range paths {
		if IsFile(path) {
			files = append(files, path)
		} else if IsDir(path) {
			dirs = append(dirs, path)
		}
	}
	return files, dirs
}

// Returns subdirs, files
func ReadDirFromPaths(paths []string, currentPath string) ([]string, []string) {
	currentPath = filepath.Clean(currentPath)
	currentPath = filepath.ToSlash(currentPath)
	var files []string
	var subdirs []string
	for _, path := range paths {
		path = filepath.ToSlash(path)
		if currentPath != "." {
			var found bool
			path, found = strings.CutPrefix(path, currentPath+"/")
			if !found {
				continue
			}
		}
		if !strings.Contains(path, "/") {
			files = append(files, path)
		} else {
			dir := strings.Split(path, "/")[0]
			if !slices.Contains(subdirs, dir) {
				subdirs = append(subdirs, dir)
			}
		}
	}
	return subdirs, files
}

// Returns all files that appear in any tree of the specified paths
func ExpandPaths(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, filepath.Clean(path))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func GetAllFilePaths() ([]string, error) {
	return ExpandPaths([]string{"."})
}

func ReadJsonStrStr(path string) (map[string]string, error) {
	stageFileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var stagingData map[string]string
	err = json.Unmarshal(stageFileData, &stagingData)
	if err != nil {
		return nil, err
	}
	return stagingData, nil
}

func WriteJsonStrStr(path string, data map[string]string) error {
	stageFileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, stageFileData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetKeys[M ~map[K]V, K comparable, V any](m M) []K {
	return slices.Collect(maps.Keys(m))
}

func SlicesDifferences[T comparable](a []T, b []T) (onlyA []T, onlyB []T) {
	aSet := make(map[T]struct{}, len(b))
	for _, value := range a {
		aSet[value] = struct{}{}
	}
	bSet := make(map[T]struct{}, len(b))
	for _, value := range a {
		bSet[value] = struct{}{}
	}
	for _, value := range a {
		if !slices.Contains(b, value) {
			onlyA = append(onlyA, value)
		}
	}
	for _, value := range b {
		if !slices.Contains(a, value) {
			onlyB = append(onlyB, value)
		}
	}
	return onlyA, onlyB
}
