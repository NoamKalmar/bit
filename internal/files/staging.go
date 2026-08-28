package files

import (
	"os"

	"github.com/noamkalmar/bit/internal/utils"
)

func CreateBlob(hash string, content []byte) error {
	file, err := os.Create(".bit/objects/blobs/" + hash)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(content)
	return nil
}

func ReadStageFile() (map[string]string, error) {
	return utils.ReadJsonStrStr(".bit/stage")
}

func WriteStageFile(stagingData map[string]string) error {
	return utils.WriteJsonStrStr(".bit/stage", stagingData)
}
