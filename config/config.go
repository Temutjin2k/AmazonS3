package config

import "errors"

var (
	Dir                  string // Storing Directory flag
	BucketMetadataFields []byte = []byte("Name,CreationDate,LastModifiedTime,Status\n")
	ObjectMetadataFields []byte = []byte("ObjectKey,Size,ContentType,LastModified\n")

	HandlerBucketList string = "BucketListHandler"
	HandlerBucket     string = "BucketHandler"
	HandlerObject     string = "ObjectHandler"

	ErrBucketExists       = errors.New("bucket already exists")
	ErrBucketDoesNotExist = errors.New("bucket does not exists")
	ErrInvalidPath        = errors.New("error deleting bucket")
	ErrNotFound           = errors.New("not Found")
)
