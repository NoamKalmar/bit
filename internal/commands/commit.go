package commands

import (
	"crypto/sha1"
	"encoding/json"
	"strings"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

func CreateCommit(message []string) error {
	commitMessage := strings.Join(message, "")
	if commitMessage == "" {
		commitMessage = "null"
	}
	commit := files.Commit{
		TreeHash:         "null",
		ParentCommitHash: "null",
		Message:          commitMessage,
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
