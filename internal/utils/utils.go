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

type PathSet map[string]struct{}

func (s PathSet) Contains(path string) bool {
	_, ok := s[filepath.Clean(path)]
	return ok
}

func (s PathSet) Add(path string) {
	s[filepath.Clean(path)] = struct{}{}
}

func (s PathSet) AddSet(paths PathSet) {
	for path := range paths {
		s.Add(path)
	}
}

func (s PathSet) Remove(path string) {
	delete(s, filepath.Clean(path))
}

func (s PathSet) Difference(other PathSet) PathSet {
	result := make(PathSet)
	for path := range s {
		if !other.Contains(path) {
			result.Add(path)
		}
	}
	return result
}

func SliceToPathSet(slice []string) PathSet {
	paths := make(PathSet)
	for _, path := range slice {
		paths.Add(path)
	}
	return paths
}

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

func CreateDirs(paths PathSet) error {
	for path := range paths {
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
func ClassifyPaths(paths PathSet) (PathSet, PathSet) {
	files := make(PathSet)
	dirs := make(PathSet)
	for path := range paths {
		if IsFile(path) {
			files.Add(path)
		} else if IsDir(path) {
			dirs.Add(path)
		}
	}
	return files, dirs
}

// Returns subdirs, files
func ReadDirFromPaths(paths PathSet, currentPath string) (PathSet, PathSet) {
	currentPath = filepath.Clean(currentPath)
	currentPath = filepath.ToSlash(currentPath)
	files := make(PathSet)
	subdirs := make(PathSet)
	for path := range paths {
		path = filepath.ToSlash(path)
		if currentPath != "." {
			var found bool
			path, found = strings.CutPrefix(path, currentPath+"/")
			if !found {
				continue
			}
		}
		if !strings.Contains(path, "/") {
			files.Add(path)
		} else {
			dir, _, _ := strings.Cut(path, "/")
			subdirs.Add(dir)
		}
	}
	return subdirs, files
}

// Returns all files that appear in any tree of the specified paths
func ExpandPaths(paths PathSet) (PathSet, error) {
	files := make(PathSet)
	for path := range paths {
		err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files.Add(path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func GetAllFilePaths() (PathSet, error) {
	dot := PathSet{}
	dot.Add(".")
	return ExpandPaths(dot)
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
