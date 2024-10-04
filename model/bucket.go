package model

import "encoding/xml"

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
	LastModified string `xml:"LastModified"`
	Status       string `xml:"Status"`
}

type BucketResponse struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}
