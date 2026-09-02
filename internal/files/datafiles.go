package files

import (
	"os"

	"github.com/noamkalmar/bit/internal/utils"
)

func ReadStageFile() (map[string]string, error) {
	return utils.ReadJsonStrStr(".bit/stage")
}

func WriteStageFile(stagingData map[string]string) error {
	return utils.WriteJsonStrStr(".bit/stage", stagingData)
}

func ReadHeadFile() (string, error) {
	content, err := os.ReadFile(".bit/head")
	return string(content), err
}

func WriteHeadFile(branchName string) error {
	return os.WriteFile(".bit/head", []byte(branchName), 0644)
}
