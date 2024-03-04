package amazon

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/flypay/go-kit/v4/internal/service"
	"github.com/flypay/go-kit/v4/pkg/objstore"
	"github.com/flypay/go-kit/v4/pkg/runtime"
)

// createLargeMessageStorage is used to bootstrap the large message storage for all environments.
func createLargeMessageStorage(sess *session.Session, runtimeCfg runtime.Config) (objstore.DownloadUploader, error) {
	if runtimeCfg.EventInMemory {
		return objstore.NewMemoryStore(), nil
	}

	if runtimeCfg.EventLargeMessageBucketLocalCreate {
		err := objstore.CreateLocalS3Buckets(sess, runtimeCfg.Environment, runtimeCfg.ServiceName, []string{service.LargeMessagesStore(runtimeCfg.ServiceName).Bucket})
		if err != nil {
			return nil, err
		}
	}

	return objstore.NewS3(sess)
}
