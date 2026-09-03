package files

import "os"

func WriteBranchRef(branchName string, hash string) error {
	return os.WriteFile(".bit/branches/"+branchName, []byte(hash), 0644)
}

func ReadBranchRef(branchName string) (string, error) {
	content, err := os.ReadFile(".bit/branches/" + branchName)
	return string(content), err
}
