package files

import (
	"os"

	"github.com/noamkalmar/bit/internal/utils"
)

func IsProjectInitialized() bool {
	return utils.IsDir(BIT_PATH)
}

func InitMainDir() error {
	// Initializing default dirs
	paths := make(utils.PathSet)
	paths.Add(BIT_PATH)
	paths.Add(BRANCHES_PATH)
	paths.Add(OBJECTS_PATH)
	paths.Add(BLOBS_PATH)
	paths.Add(TREES_PATH)
	paths.Add(COMMITS_PATH)
	err := utils.CreateDirs(paths)
	if err != nil {
		return err
	}
	// Initializing default files
	err = os.WriteFile(DEFAULT_BRANCH_PATH, []byte{}, 0644)
	err = os.WriteFile(INDEX_PATH, []byte("{}"), 0644)
	// Head file should be initialized to point to the main branch
	err = os.WriteFile(HEAD_PATH, []byte(DEFAULT_BRANCH), 0644)
	return err
}
