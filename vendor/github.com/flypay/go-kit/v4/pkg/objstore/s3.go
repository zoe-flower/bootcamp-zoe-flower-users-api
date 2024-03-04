package objstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/pkg/errors"
)

const defaultRegion = "eu-west-1"

var (
	// ErrNoSuchKey is returned when a key is not found in S3. It's the S3 equivalent of projections.ErrNotFound
	ErrNoSuchKey = errors.New("no such object key")

	// ErrNoSuchBucket is returned when the method tries to access a bucket that does not exist.
	// This error mainly exists for the copy object method.
	ErrNoSuchBucket = errors.New("no such bucket")
)

type S3 struct {
	downloader s3Downloader
	uploader   s3Uploader
	deleter    s3ObjectDeleter
	getter     s3ObjectGetter
	copier     s3ObjectCopier
	header     s3Header
	metrics    s3Metrics
}

// NewS3 returns an AWS S3 Download and Upload manager
func NewS3(cfg *session.Session) (*S3, error) {
	if cfg == nil {
		var err error
		cfg, err = session.NewSession(aws.NewConfig().WithRegion(defaultRegion))
		if err != nil {
			return nil, err
		}
	}

	s3client := awss3.New(cfg)
	xray.AWS(s3client.Client)

	s3metrics := RegisterS3Metrics()

	return &S3{
		downloader: s3manager.NewDownloaderWithClient(s3client),
		uploader:   s3manager.NewUploaderWithClient(s3client),
		deleter:    s3client,
		copier:     s3client,
		getter:     s3client,
		header:     s3client,
		metrics:    s3metrics,
	}, nil
}

func (s S3) ObjectMetadata(ctx context.Context, bucket Bucket, key string) (*Metadata, error) {
	start := time.Now()

	resp, err := s.header.HeadObjectWithContext(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(string(bucket)),
		Key:    aws.String(key),
	})

	s.metrics.MetricsReq(start, "ObjectMetadata", err)

	if err != nil {
		return nil, s.wrapAwsError(err)
	}
	return metadataFromAWSHead(resp), nil
}

func (s S3) ObjectMetadataURL(ctx context.Context, objectURL *url.URL) (*Metadata, error) {
	start := time.Now()

	bucket, key := parseBucketAndKeyFromURL(*objectURL)
	resp, err := s.header.HeadObjectWithContext(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	s.metrics.MetricsReq(start, "ObjectMetadataURL", err)

	if err != nil {
		return nil, s.wrapAwsError(err)
	}
	return metadataFromAWSHead(resp), nil
}

func (s S3) GetObject(ctx context.Context, bkt Bucket, key string) (io.ReadCloser, error) {
	start := time.Now()

	resp, err := s.getter.GetObjectWithContext(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(string(bkt)),
		Key:    aws.String(key),
	})
	s.metrics.MetricsReq(start, "GetObject", err)

	return resp.Body, s.wrapAwsError(err)
}

func (s S3) GetObjectURL(ctx context.Context, objectURL *url.URL) (io.ReadCloser, error) {
	start := time.Now()

	objInput := buildObjectInputFromURL(*objectURL)
	resp, err := s.getter.GetObjectWithContext(ctx, &objInput)
	s.metrics.MetricsReq(start, "GetObjectURL", err)

	return resp.Body, s.wrapAwsError(err)
}

func (s S3) ListObjects(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error) {
	start := time.Now()

	input := &awss3.ListObjectsV2Input{
		Bucket: aws.String(string(bucket)),
	}

	if directoryPath != "" {
		input.Prefix = aws.String(directoryPath)
	}

	keys, err := s.addObjectKeys(ctx, input, []string{})
	s.metrics.MetricsReq(start, "ListObjects", err)

	return keys, err
}

func (s S3) ListSubdirectories(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error) {
	start := time.Now()

	input := &awss3.ListObjectsV2Input{
		Bucket:    aws.String(string(bucket)),
		Delimiter: aws.String("/"),
	}

	if directoryPath != "" {
		input.Prefix = aws.String(directoryPath)
	}

	keys, err := s.addSubdirectories(ctx, input, []string{})
	s.metrics.MetricsReq(start, "ListSubdirectories", err)

	return keys, err
}

func (s S3) Download(ctx context.Context, bucket Bucket, key string, body io.WriterAt) (written int64, err error) {
	start := time.Now()

	written, err = s.downloader.DownloadWithContext(ctx, body, &awss3.GetObjectInput{
		Bucket: aws.String(string(bucket)),
		Key:    aws.String(key),
	})

	s.metrics.MetricsReq(start, "Download", err)

	return written, s.wrapAwsError(err)
}

func (s S3) DownloadObject(ctx context.Context, bucket Bucket, key string) ([]byte, error) {
	var err error
	start := time.Now()
	defer s.metrics.MetricsReq(start, "DownloadObject", err)

	resp, err := s.getter.GetObjectWithContext(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(string(bucket)),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, s.wrapAwsError(err)
	}

	var buf bytes.Buffer
	// If there is a content length set in the response headers, use that
	// to pre-initialise a buffer of that size.
	if resp.ContentLength != nil {
		buf.Grow(int(*resp.ContentLength))
	}
	_, err = buf.ReadFrom(resp.Body)

	return buf.Bytes(), err
}

func (s S3) DownloadURL(ctx context.Context, objectURL *url.URL, body io.WriterAt) (written int64, err error) {
	start := time.Now()
	defer s.metrics.MetricsReq(start, "DownloadURL", err)
	objInput := buildObjectInputFromURL(*objectURL)
	written, err = s.downloader.DownloadWithContext(ctx, body, &objInput)
	return written, s.wrapAwsError(err)
}

func (s S3) DownloadObjectURL(ctx context.Context, objectURL *url.URL) ([]byte, error) {
	var err error
	start := time.Now()
	defer s.metrics.MetricsReq(start, "DownloadObjectURL", err)

	bucket, key := parseBucketAndKeyFromURL(*objectURL)
	resp, err := s.getter.GetObjectWithContext(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, s.wrapAwsError(err)
	}

	var buf bytes.Buffer
	// If there is a content length set in the response headers, use that
	// to pre-initialise a buffer of that size.
	if resp.ContentLength != nil {
		buf.Grow(int(*resp.ContentLength))
	}
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}

func (s S3) Upload(ctx context.Context, bucket Bucket, key string, body io.Reader) (location string, err error) {
	start := time.Now()
	defer s.metrics.MetricsReq(start, "Upload", err)

	output, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(string(bucket)),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}

	return output.Location, err
}

func (s S3) DeleteObject(ctx context.Context, bucket Bucket, key string) (err error) {
	start := time.Now()
	defer s.metrics.MetricsReq(start, "DeleteObject", err)

	_, err = s.deleter.DeleteObjectWithContext(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(string(bucket)),
		Key:    aws.String(key),
	})
	return err
}

func (s S3) CopyObject(ctx context.Context, srcLocation string, destBucket Bucket, dest string) (err error) {
	start := time.Now()
	defer s.metrics.MetricsReq(start, "CopyObject", err)

	var u *url.URL
	u, err = url.Parse(srcLocation)
	if err != nil {
		return errors.Wrapf(err, "parsing src location %q", srcLocation)
	}
	err = s.copyObjectURL(ctx, u, destBucket, dest)
	return err
}

func (s S3) CopyObjectURL(ctx context.Context, objectURL *url.URL, destBucket Bucket, dest string) (err error) {
	start := time.Now()
	defer s.metrics.MetricsReq(start, "CopyObjectURL", err)

	err = s.copyObjectURL(ctx, objectURL, destBucket, dest)
	return err
}

func (s S3) copyObjectURL(ctx context.Context, objectURL *url.URL, destBucket Bucket, dest string) (err error) {
	srcBucket, srcKey := parseBucketAndKeyFromURL(*objectURL)
	fullpath := fmt.Sprintf("%s/%s", srcBucket, srcKey)
	_, err = s.copier.CopyObjectWithContext(ctx, &awss3.CopyObjectInput{
		Bucket:     aws.String(string(destBucket)),
		Key:        aws.String(dest),
		CopySource: aws.String(copyQueryEscape(fullpath)),
	})
	return s.wrapAwsError(err)
}

// copyQueryEscape will query escape the copy path unless the runtime configuration
// says not to. This is an issue with Minio which requires the copy source path
// to NOT be url encoded, but the S3 api expects the copy source path to be
// url encoded.
// Minio have insisted in a few issues that this is working as expected:
// - https://github.com/minio/minio/issues/7277
func copyQueryEscape(path string) string {
	if cfg := runtime.DefaultConfig(); cfg.ObjstoreNoCopyQueryEscape {
		return path
	}
	return url.QueryEscape(path)
}

// MakeEnvAwareSession creates a session that will adapt to different environments
func MakeEnvAwareSession() (*session.Session, error) {
	return MakeEnvAwareSessionWithRuntimeConfig(runtime.DefaultConfig())
}

func MakeEnvAwareSessionWithRuntimeConfig(cfg runtime.Config) (*session.Session, error) {
	if cfg.ObjstoreEndpoint != "" {
		return newCustomSession(cfg.ObjstoreEndpoint)
	}

	return newDefaultSession()
}

func newDefaultSession() (*session.Session, error) {
	return session.NewSession(aws.NewConfig().WithRegion(defaultRegion))
}

// newCustomSession creates an AWS session for use on a developer machine with Minio
func newCustomSession(endpoint string) (*session.Session, error) {
	return session.NewSession(
		aws.NewConfig().
			WithCredentialsChainVerboseErrors(true).
			WithCredentials(credentials.NewStaticCredentials("flyt", "flyt1234", "")).
			WithEndpoint(endpoint).
			WithS3ForcePathStyle(true).
			WithDisableSSL(true).
			WithRegion("eu-west-1"),
	)
}

// parseBucketAndKeyFromURL will parse the URIs to get the bucket and key
func parseBucketAndKeyFromURL(url url.URL) (string, string) {
	// When using Minio the bucket name is located differently
	if strings.Contains(url.Host, "localhost") {
		reqpath := url.Path
		bucket := strings.Split(reqpath, "/")[1]
		key := strings.Join(strings.Split(reqpath, "/")[2:], "/")
		return bucket, key
	}

	hostParts := strings.SplitN(url.Host, ".", 2)
	bucket := hostParts[0]
	key := url.Path
	if strings.HasPrefix(url.Path, "/") {
		key = key[1:]
	}
	return bucket, key
}

// buildObjextInputFromURL is used to parse URIs for convenience functions GetObjectURL and DownloadURL
func buildObjectInputFromURL(url url.URL) awss3.GetObjectInput {
	bucket, key := parseBucketAndKeyFromURL(url)
	return awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
}

func (s S3) addObjectKeys(ctx context.Context, input *awss3.ListObjectsV2Input, keys []string) ([]string, error) {
	resp, err := s.getter.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return keys, s.wrapAwsError(err)
	}

	for _, c := range resp.Contents {
		keys = append(keys, *c.Key)
	}

	// Response will be truncated if len(resp.Contents) > 1000
	if resp.IsTruncated != nil {
		if *resp.IsTruncated {
			input.ContinuationToken = resp.NextContinuationToken
			keys, err = s.addObjectKeys(ctx, input, keys)
			if err != nil {
				return keys, err
			}
		}
	}

	return keys, nil
}

func (s S3) addSubdirectories(ctx context.Context, input *awss3.ListObjectsV2Input, keys []string) ([]string, error) {
	resp, err := s.getter.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return keys, s.wrapAwsError(err)
	}

	for _, c := range resp.CommonPrefixes {
		keys = append(keys, *c.Prefix)
	}

	// Response will be truncated if len(resp.CommonPrefixes) > 1000
	if resp.IsTruncated != nil && *resp.IsTruncated {
		input.ContinuationToken = resp.NextContinuationToken
		keys, err = s.addSubdirectories(ctx, input, keys)
		if err != nil {
			return keys, err
		}
	}

	return keys, nil
}

// wrapAwsError gives a more clean way to match for specific S3 specific errors, and wraps the underlying error from the aws sdk.
// i.e you could read a restaurant key in an object store and match the returned error with:
/*
	if errors.Is(err, objstore.ErrNoSuchKey) {

	}
*/
// while not have to deal with matching errors based on string value.
// Note: The method does not work for the v2 sdk.
func (s S3) wrapAwsError(err error) error {
	var s3err awserr.Error

	if ok := errors.As(err, &s3err); ok {
		switch s3err.Code() {
		case awss3.ErrCodeNoSuchKey:
			return fmt.Errorf("%v: %w", err, ErrNoSuchKey)
		case awss3.ErrCodeNoSuchBucket:
			return fmt.Errorf("%v: %w", err, ErrNoSuchBucket)
		}
	}
	return err
}

type s3Uploader interface {
	UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type s3Downloader interface {
	DownloadWithContext(aws.Context, io.WriterAt, *awss3.GetObjectInput, ...func(*s3manager.Downloader)) (n int64, err error)
}

type s3ObjectDeleter interface {
	DeleteObjectWithContext(aws.Context, *awss3.DeleteObjectInput, ...request.Option) (*awss3.DeleteObjectOutput, error)
}

type s3ObjectCopier interface {
	CopyObjectWithContext(aws.Context, *awss3.CopyObjectInput, ...request.Option) (*awss3.CopyObjectOutput, error)
}

type s3ObjectGetter interface {
	GetObjectWithContext(context.Context, *awss3.GetObjectInput, ...request.Option) (*awss3.GetObjectOutput, error)
	ListObjectsV2WithContext(ctx context.Context, input *awss3.ListObjectsV2Input, options ...request.Option) (*awss3.ListObjectsV2Output, error)
}

type s3Header interface {
	HeadObjectWithContext(context.Context, *awss3.HeadObjectInput, ...request.Option) (*awss3.HeadObjectOutput, error)
}

type s3Metrics interface {
	MetricsReq(startTime time.Time, reqType string, err error)
}
