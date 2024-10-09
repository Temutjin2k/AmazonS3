package utils

import (
	"os"

	"triple-s/config"
)

func MakeDir(path string) error {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}

	// Creating buckets.csv to store metadata of buckets
	err = os.WriteFile(path+"/buckets.csv", config.BucketMetadataFields, 0o755)
	if err != nil {
		return err
	}
	return nil
}
