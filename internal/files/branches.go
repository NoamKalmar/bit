package files

import (
	"os"
	"path/filepath"
)

func WriteBranchRef(branchName string, hash string) error {
	return os.WriteFile(filepath.Join(BRANCHES_PATH, branchName), []byte(hash), 0644)
}

func ReadBranchRef(branchName string) (string, error) {
	content, err := os.ReadFile(filepath.Join(BRANCHES_PATH, branchName))
	return string(content), err
}
