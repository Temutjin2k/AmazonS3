package utils

import (
	"triple-s/config"
)

func IsBucketExist(bucket string) (bool, error) {
	isBucketExist, err := SearchValueCSV(config.Dir+"/buckets.csv", "Name", bucket)
	if err != nil {
		return false, err
	}
	if !isBucketExist {
		return false, config.ErrBucketExists
	}
	return true, nil
}

func IsObjectExist(bucket, objectKey string) (bool, error) {
	isObjectKeyExist, err := SearchValueCSV(config.Dir+"/"+bucket+"/objects.csv", "ObjectKey", objectKey)
	if err != nil {
		return false, err
	}
	if !isObjectKeyExist {
		return false, config.ErrObjectDoesNotExist
	}
	return true, nil
}
