package files

import "os"

func CreateBlob(hash string, content []byte) error {
	file, err := os.Create(".bit/objects/blobs/" + hash)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(content)
	return nil
}
