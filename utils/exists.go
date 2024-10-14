package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"triple-s/config"
)

func IsBucketExist(bucket string) bool {
	isBucketExist, err := SearchValueCSV(filepath.Join(config.Dir, "buckets.csv"), "Name", bucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return isBucketExist
}

func IsObjectExist(bucket, objectKey string) bool {
	isObjectKeyExist, err := SearchValueCSV(filepath.Join(config.Dir, bucket, "objects.csv"), "ObjectKey", objectKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return isObjectKeyExist
}
