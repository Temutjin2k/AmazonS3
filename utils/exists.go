package utils

import (
	"path/filepath"

	"triple-s/config"
)

func IsBucketExist(bucket string) (bool, error) {
	isBucketExist, err := SearchValueCSV(filepath.Join(config.Dir, "buckets.csv"), "Name", bucket)
	if err != nil {
		return false, err
	}
	if !isBucketExist {
		return false, nil
	}
	return true, nil
}

func IsObjectExist(bucket, objectKey string) (bool, error) {
	isObjectKeyExist, err := SearchValueCSV(filepath.Join(config.Dir, bucket, "objects.csv"), "ObjectKey", objectKey)
	if err != nil {
		return false, err
	}
	if !isObjectKeyExist {
		return false, nil
	}
	return true, nil
}
