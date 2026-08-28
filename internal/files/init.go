package files

import (
	"os"

	"github.com/noamkalmar/bit/internal/utils"
)

func IsProjectInitalized() bool {
	info, err := os.Stat(".bit")
	if err == nil {
		return info.IsDir()
	}
	return false
}

func InitMainDir() error {
	// Initalizing default dirs
	paths := []string{
		".bit/objects/commits",
		".bit/objects/trees",
		".bit/objects/blobs",
		".bit/branches",
	}
	err := utils.CreateDirs(paths)
	if err != nil {
		return err
	}
	// Initalizing default files
	err = os.WriteFile(".bit/branches/main", []byte{}, 0644)
	err = os.WriteFile(".bit/stage", []byte("{}"), 0644)
	// Head file should be initalized to point to the main branch
	err = os.WriteFile(".bit/head", []byte("main"), 0644)
	return err
}
