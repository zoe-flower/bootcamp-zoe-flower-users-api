package objstore

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/flypay/go-kit/v4/pkg/runtime"
)

// Bucket is a type that defines where to store objects
type Bucket string

// GetServiceBucketName creates a bucket name based on the service and environment
func GetServiceBucketName(bucketName, serviceName string) Bucket {
	cfg := runtime.DefaultConfig()
	sName := cfg.ServiceName
	if sName == "" {
		sName = serviceName
	}
	return Bucket(serviceBucketName(cfg.Environment, sName, bucketName))
}

func serviceBucketName(environment, service, bucket string) string {
	return strings.Join([]string{environment, service, bucket}, "-")
}

// ResolveObjectStorage returns an implementation of whatever storage we need
func ResolveObjectStorage(session *session.Session) (DownloadUploader, error) {
	return ResolveObjectStorageWithRuntimeConfig(session, runtime.DefaultConfig())
}

func ResolveObjectStorageWithRuntimeConfig(session *session.Session, cfg runtime.Config) (DownloadUploader, error) {
	if cfg.ObjstoreInMemory {
		return NewMemoryStore(), nil
	}

	if cfg.ObjstoreLocalCreate {
		err := CreateLocalS3Buckets(session, cfg.Environment, cfg.ServiceName, cfg.ObjstoreBuckets)
		if err != nil {
			return nil, err
		}
	}

	return NewS3(session)
}

// Downloader can be used to fetch an object from the object storage. It uses concurrent download
// to increase bandwidth.
// Deprecated: use ObjectDownloader instead
type Downloader interface {
	Download(ctx context.Context, bucket Bucket, key string, body io.WriterAt) (written int64, err error)
}

// URLDownloader is similar to Downloader with a url.URL parameter
// Deprecated: use URLObjectDownloader instead
type URLDownloader interface {
	DownloadURL(ctx context.Context, objectURL *url.URL, body io.WriterAt) (written int64, err error)
}

// Uploader puts the data read from the reader into the bucket at the given key
type Uploader interface {
	Upload(ctx context.Context, bucket Bucket, key string, body io.Reader) (location string, err error)
}

// ObjectGetter implements the fetching of an object and returns a reader so that you can stream the data into your application
type ObjectGetter interface {
	GetObject(ctx context.Context, bkt Bucket, key string) (io.ReadCloser, error)
}

// URLObjectGetter is similar to ObjectGetter with a url.URL parameter
type URLObjectGetter interface {
	GetObjectURL(ctx context.Context, objectURL *url.URL) (io.ReadCloser, error)
}

// ObjectLister returns keys of all objects within a bucket
// directoryPath is an optional argument to limit results to a specific directory, pass as empty string for all objects in the bucket
type ObjectLister interface {
	ListObjects(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error)
}

// SubdirectoriesLister returns a list of all subdirectories in the provided path
// directoryPath is an optional argument to limit results to a specific directory
type SubdirectoriesLister interface {
	ListSubdirectories(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error)
}

// ObjectMetadata will perform a HEAD request on the object in the bucket
// and return metadata
type ObjectMetadata interface {
	ObjectMetadata(ctx context.Context, bucket Bucket, key string) (*Metadata, error)
}

// URLObjectMetadata will perform a HEAD request on the object in the bucket
// and return metadata using a url.URL
type URLObjectMetadata interface {
	ObjectMetadataURL(ctx context.Context, objectURL *url.URL) (*Metadata, error)
}

// ObjectDownloader will download an object into memory and return the byte slice
type ObjectDownloader interface {
	DownloadObject(ctx context.Context, bucket Bucket, key string) (data []byte, err error)
}

// URLObjectDownloader will download an object into memory from a url.URL and return the byte slice
type URLObjectDownloader interface {
	DownloadObjectURL(ctx context.Context, objectURL *url.URL) (data []byte, err error)
}

// ObjectCopier will copy your object from full `location` (incl. bucket name) to destination (bucket, key).
// The location will be url encoded.
type ObjectCopier interface {
	CopyObject(ctx context.Context, location string, bucket Bucket, key string) error
}

// ObjectCopier will copy your object by a src url.URL
type URLObjectCopier interface {
	CopyObjectURL(ctx context.Context, objectURL *url.URL, bucket Bucket, dest string) error
}

// ObjectDeleter deletes an object from the specified bucket and key
type ObjectDeleter interface {
	DeleteObject(ctx context.Context, bucket Bucket, key string) error
}

// DownloadUploader encapulsates all the things you can do with the objstore client
type DownloadUploader interface {
	// Deprecated: use ObjectDownloader instead
	Downloader
	// Deprecated: use URLObjectDownloader instead
	URLDownloader

	// Retrieving objects
	ObjectGetter
	URLObjectGetter
	ObjectDownloader
	URLObjectDownloader

	// Listing objects
	ObjectLister
	// Listing subdirectories
	SubdirectoriesLister
	// Storing objects
	Uploader

	// Metadata for objects
	ObjectMetadata
	URLObjectMetadata

	// Copy objects
	ObjectCopier
	URLObjectCopier

	// Delete objects
	ObjectDeleter
}

func CreateLocalS3Buckets(session *session.Session, environment, serviceName string, buckets []string) error {
	s3Client := awsS3.New(session)

	for _, bucket := range buckets {
		bucketName := serviceBucketName(environment, serviceName, bucket)

		b := &awsS3.CreateBucketInput{
			Bucket: &bucketName,
		}

		_, err := s3Client.CreateBucket(b)

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case awsS3.ErrCodeBucketAlreadyOwnedByYou:
				fallthrough
			case awsS3.ErrCodeBucketAlreadyExists:
				continue
			default:
				return aerr
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
