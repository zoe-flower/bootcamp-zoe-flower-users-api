package objstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
)

type memUpload struct {
	Body         []byte
	URL          url.URL
	LastModified time.Time
}

type MemoryStore struct {
	buckets map[string]memUpload
	mu      sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets: make(map[string]memUpload),
	}
}

func makeFakeURL(bucket Bucket, key string) url.URL {
	return url.URL{
		Scheme: "http",
		Host:   "example.com",
		Path:   fmt.Sprintf("%s/%s", bucket, key),
	}
}

func makeBucketIdent(bucket Bucket, key string) string {
	return fmt.Sprintf("%s/%s", bucket, key)
}

func splitBucketKey(fullkey string) (bucket Bucket, key string, ok bool) {
	idx := strings.Index(fullkey, "/")
	if idx <= -1 {
		return "", "", false
	}

	bucket = Bucket(fullkey[:idx])
	key = fullkey[idx+1:]
	ok = true
	return
}

func (ms *MemoryStore) Upload(ctx context.Context, bucket Bucket, key string, body io.Reader) (location string, err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var b []byte
	b, err = io.ReadAll(body)
	if err != nil {
		return
	}

	upload := memUpload{
		Body:         b,
		URL:          makeFakeURL(bucket, key),
		LastModified: time.Now(),
	}

	ms.buckets[makeBucketIdent(bucket, key)] = upload

	location = upload.URL.String()

	return
}

func (ms *MemoryStore) get(bucket Bucket, key string) (memUpload, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	upload, ok := ms.buckets[makeBucketIdent(bucket, key)]
	return upload, ok
}

func (ms *MemoryStore) getURL(url *url.URL) (memUpload, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, upload := range ms.buckets {
		if upload.URL.String() == url.String() {
			return upload, true
		}
	}
	return memUpload{}, false
}

func (ms *MemoryStore) Download(ctx context.Context, bucket Bucket, key string, body io.WriterAt) (int64, error) {
	upload, ok := ms.get(bucket, key)
	if ok {
		wrote, err := body.WriteAt(upload.Body, 0)
		return int64(wrote), err
	}

	return 0, ErrNoSuchKey
}

func (ms *MemoryStore) DownloadURL(ctx context.Context, url *url.URL, body io.WriterAt) (int64, error) {
	upload, ok := ms.getURL(url)
	if ok {
		wrote, err := body.WriteAt(upload.Body, 0)
		return int64(wrote), err
	}

	return 0, ErrNoSuchKey
}

func (ms *MemoryStore) GetObject(ctx context.Context, bucket Bucket, key string) (io.ReadCloser, error) {
	upload, ok := ms.get(bucket, key)
	if ok {
		return io.NopCloser(bytes.NewReader(upload.Body)), nil
	}

	return nil, ErrNoSuchKey
}

func (ms *MemoryStore) GetObjectURL(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	upload, ok := ms.getURL(url)

	if ok {
		return io.NopCloser(bytes.NewReader(upload.Body)), nil
	}

	return nil, ErrNoSuchKey
}

func (ms *MemoryStore) ListObjects(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var objKeys []string

	prefix := makeBucketIdent(bucket, directoryPath)
	for b := range ms.buckets {
		if strings.HasPrefix(b, prefix) {
			objKeys = append(objKeys, removeBucketFromPath(b, bucket))
		}
	}

	return objKeys, nil
}

func (ms *MemoryStore) ListSubdirectories(ctx context.Context, bucket Bucket, directoryPath string) ([]string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var subdirectories []string

	prefix := makeBucketIdent(bucket, directoryPath)
	for b := range ms.buckets {
		if strings.HasPrefix(b, prefix) {
			objects := strings.SplitAfter(b, prefix)[1]
			subdir := strings.SplitAfter(objects, "/")[0]
			if subdir[len(subdir)-1:] == "/" {
				subdirectories = append(subdirectories, directoryPath+subdir)
			}
		}
	}

	return lo.Uniq(subdirectories), nil
}

func (ms *MemoryStore) DownloadObject(ctx context.Context, bucket Bucket, key string) ([]byte, error) {
	upload, _ := ms.get(bucket, key)
	return upload.Body, nil
}

func (ms *MemoryStore) DownloadObjectURL(ctx context.Context, url *url.URL) ([]byte, error) {
	upload, _ := ms.getURL(url)
	return upload.Body, nil
}

func (ms *MemoryStore) ObjectMetadata(ctx context.Context, bucket Bucket, key string) (*Metadata, error) {
	upload, _ := ms.get(bucket, key)
	md := Metadata{
		ContentLength: int64(len(upload.Body)),
		LastModified:  upload.LastModified,
	}
	return &md, nil
}

func (ms *MemoryStore) ObjectMetadataURL(ctx context.Context, url *url.URL) (*Metadata, error) {
	upload, _ := ms.getURL(url)
	md := Metadata{
		ContentLength: int64(len(upload.Body)),
		LastModified:  upload.LastModified,
	}
	return &md, nil
}

func (ms *MemoryStore) CopyObject(ctx context.Context, srcLocation string, destBucket Bucket, destKey string) error {
	bucket, key, ok := splitBucketKey(srcLocation)
	if !ok {
		return errors.New("Bad existing bucket/key reference")
	}
	m, _ := ms.get(bucket, key)
	b := make([]byte, len(m.Body))
	copy(b, m.Body)
	if _, err := ms.Upload(ctx, destBucket, destKey, bytes.NewReader(m.Body)); err != nil {
		return err
	}
	return nil
}

func (ms *MemoryStore) CopyObjectURL(ctx context.Context, url *url.URL, bucket Bucket, key string) error {
	m, _ := ms.getURL(url)
	b := make([]byte, len(m.Body))
	copy(b, m.Body)
	if _, err := ms.Upload(ctx, bucket, key, bytes.NewReader(b)); err != nil {
		return err
	}
	return nil
}

func (ms *MemoryStore) DeleteObject(ctx context.Context, bucket Bucket, key string) error {
	if _, ok := ms.get(bucket, key); ok {
		ms.mu.Lock()
		defer ms.mu.Unlock()
		delete(ms.buckets, makeBucketIdent(bucket, key))
	}
	return nil
}

func removeBucketFromPath(path string, bucket Bucket) string {
	return strings.TrimPrefix(path, fmt.Sprintf("%v/", bucket))
}

type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}
