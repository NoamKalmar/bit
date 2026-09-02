package files

import "os"

func WriteBranchRef(branchName string, hash string) error {
	return os.WriteFile(".bit/branches/"+branchName, []byte(hash), 0644)
}
