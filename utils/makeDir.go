package utils

import (
	"os"

	"triple-s/config"
)

func MakeDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return err
	}

	// Creating buckets.csv to store metadata of bucket
	err = os.WriteFile(path+"/buckets.csv", config.BucketMetadataFields, 0o755)
	if err != nil {
		return err
	}
	return nil
}
