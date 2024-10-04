package model

import "encoding/xml"

type Bucket struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
	LastModified string `xml:"LastModified"`
}

type BucketResponse struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}
