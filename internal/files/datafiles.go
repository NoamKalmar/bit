package files

import (
	"os"

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
