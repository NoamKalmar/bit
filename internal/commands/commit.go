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
	branch, err := files.ReadHeadFile()
	if err != nil {
		return err
	}
	parent, err := files.ReadBranchRef(branch) // last commit is the parent of the new commit
	if err != nil {
		return err
	}
	commit := files.Commit{
		TreeHash:         "",
		ParentCommitHash: parent,
		Message:          commitMessage,
	}
	content, err := json.MarshalIndent(commit, "", "    ")
	if err != nil {
		return err
	}
	hasher := sha1.New()
	hash := utils.GetHash(content, hasher)
	files.CreateCommit(hash, content) // Create a new commit object
	files.WriteBranchRef(branch, hash)
	return nil
}
