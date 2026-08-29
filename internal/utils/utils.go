package utils

import (
	"encoding/hex"
	"encoding/json"
	"hash"
	"os"
	"path/filepath"
)

func GetHash(content []byte, hasher hash.Hash) string {
	hasher.Write(content)
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
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

func ExpandPaths(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, nil
		}
	}
	return files, nil
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

func WriteJsonStrStr(path string, stagingData map[string]string) error {
	stageFileData, err := json.MarshalIndent(stagingData, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, stageFileData, 0644)
	if err != nil {
		return err
	}
	return nil
}
