package config

var (
	Dir                  string // Storing Directory flag
	BucketMetadataFields []byte = []byte("Name,CreationDate,LastModifiedTime,Status\n")
	ObjectMetadataFields []byte = []byte("ObjectKey,Size,ContentType,LastModified\n")
)
