package commands

import (
	"crypto/sha1"
	"encoding/json"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

func CreateCommit(commitMessage string) error {
	commit := files.Commit{
		TreeHash:         "null",
		ParentCommitHash: "null",
	}
	content, err := json.MarshalIndent(commit, "", "    ")
	if err != nil {
		return err
	}
	hasher := sha1.New()
	hash := utils.GetHash(content, hasher)
	files.CreateCommit(hash, content)
	branch, err := files.ReadHeadFile()
	if err != nil {
		return err
	}
	files.WriteBranchRef(branch, hash)
	return nil
}
