package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"hash"
	"os"
)

func GetHash(content []byte, hasher hash.Hash) string {
	hasher.Write(content)
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func fileToBlob(path string, hasher hash.Hash) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	// Create a blob with the file content's hash as the name
	hash := GetHash(content, hasher)
	file, err := os.Create(".bit/objects/blobs/" + hash)
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.WriteString(string(content))
	return hash, nil
}

func ReadStageFile() (map[string]string, error) {
	stageFileData, err := os.ReadFile(".bit/stage")
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

func WriteStageFile(stagingData map[string]string) error {
	stageFileData, err := json.MarshalIndent(stagingData, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(".bit/stage", stageFileData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func StageFiles(paths []string) error {
	stageData, err := ReadStageFile()
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
	WriteStageFile(stageData)
	return nil
}
