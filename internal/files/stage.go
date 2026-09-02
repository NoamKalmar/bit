package files

import (
	"github.com/noamkalmar/bit/internal/utils"
)

func ReadStageFile() (map[string]string, error) {
	return utils.ReadJsonStrStr(".bit/stage")
}

func WriteStageFile(stagingData map[string]string) error {
	return utils.WriteJsonStrStr(".bit/stage", stagingData)
}
