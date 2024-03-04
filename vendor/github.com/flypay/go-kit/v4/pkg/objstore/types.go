package objstore

import (
	"time"

	awss3 "github.com/aws/aws-sdk-go/service/s3"
)

// Metadata contains metadata from aws
type Metadata struct {
	awsHeadObjectOutput *awss3.HeadObjectOutput
	ContentLength       int64
	LastModified        time.Time
	ContentType         string
}

func metadataFromAWSHead(headObject *awss3.HeadObjectOutput) *Metadata {
	var contentLength int64
	if headObject.ContentLength != nil {
		contentLength = *headObject.ContentLength
	}
	var contentType string
	if headObject.ContentType != nil {
		contentType = *headObject.ContentType
	}
	var lastModified time.Time
	if headObject.LastModified != nil {
		lastModified = *headObject.LastModified
	}
	return &Metadata{
		ContentLength:       contentLength,
		ContentType:         contentType,
		LastModified:        lastModified,
		awsHeadObjectOutput: headObject,
	}
}

func (m *Metadata) AWSHeadObjectOutput() *awss3.HeadObjectOutput {
	return m.awsHeadObjectOutput
}
