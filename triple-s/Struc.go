package main

import (
	"encoding/xml"
)

// Bucket represents the metadata of a bucket.
type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"` // Added status field
}

// BucketList represents a list of buckets.
type BucketList struct {
	XMLName xml.Name `xml:"Buckets"`
	Buckets []Bucket `xml:"Bucket"`
}

type Object struct {
	ObjectKey    string `xml:"ObjectKey"`
	Size         int64  `xml:"Size"`
	ContentType  string `xml:"ContentType"`
	LastModified string `xml:"LastModified"`
}

// ObjectList represents a list of objects.
type ObjectList struct {
	XMLName xml.Name `xml:"Objects"`
	Objects []Object `xml:"Object"`
}

type ErrorResponse struct {
	Code    int      `xml:"Code"`
	Message string   `xml:"Message"`
	XMLName xml.Name `xml:"Error"`
}
