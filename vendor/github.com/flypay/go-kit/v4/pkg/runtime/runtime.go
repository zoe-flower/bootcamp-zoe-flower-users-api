package runtime

import (
	"fmt"
	"sync"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/flypay/go-kit/v4/internal/service"
)

var (
	cfg Config
	mu  sync.Mutex
)

func init() {
	_, _ = ResolveConfigAndSetDefault(WithFromENV())
}

// DefaultConfig will return the config set up by 'ResolveConfigAndSetDefault'
// This is usually set up once at start up
func DefaultConfig() Config {
	mu.Lock()
	defer mu.Unlock()
	return cfg
}

type InterfaceType struct {
	Capability interface{}
}

// Config provides a struct of all the environment variables used throughout go-kit
type Config struct {
	Environment   string `env:"CONFIG_APP_ENVIRONMENT"`
	ServiceName   string `env:"CONFIG_APP_NAME"`
	Team          string `env:"CONFIG_APP_TEAM"`
	Flow          string `env:"CONFIG_APP_FLOW"`
	Tier          string `env:"CONFIG_APP_TIER" envDefault:"Undefined"`
	Version       string `env:"CONFIG_APP_VERSION" envDefault:"development"`
	Production    bool   `env:"CONFIG_APP_IN_PROD"`
	UseGha        bool   `env:"CONFIG_APP_USE_GHA"`
	IsGhe         bool   `env:"CONFIG_APP_IS_GHE"`
	DisableSentry bool   `env:"CONFIG_APP_DISABLE_SENTRY"`

	SentryDSN       string `env:"CONFIG_SENTRY_DSN"`
	XrayDaemon      string `env:"CONFIG_XRAY_DAEMON"`
	XrayLocalLogger bool   `env:"CONFIG_XRAY_LOCAL_LOGGER"`

	PanicOnRecover bool `env:"CONFIG_PANIC_ON_RECOVER"`

	AWSAccountID string `env:"CONFIG_AWS_ACCOUNT_ID"`
	AWSRegion    string `env:"AWS_REGION" envDefault:"eu-west-1"`

	HTTPPort           string        `env:"APP_CONFIG_HTTP_PORT" envDefault:"8080"`
	HTTPHandlerTimeout time.Duration `env:"CONFIG_APP_HTTP_HANDLER_TIMEOUT" envDefault:"5s"`

	EventWorkerPool                    int      `env:"CONFIG_APP_EVENT_WORKER_POOL" envDefault:"25"`
	EventHandlerTimeout                string   `env:"CONFIG_APP_EVENT_HANDLER_TIMEOUT" envDefault:"30s"`
	EventHandlerVisibilityTimeout      string   `env:"CONFIG_APP_EVENT_VISIBILITY_TIMEOUT" envDefault:"35s"`
	EventSQSTopicName                  string   `env:"CONFIG_EVENT_BUS_SQS_TOPIC_NAME"`
	EventSQSQueueName                  string   `env:"CONFIG_EVENT_BUS_SQS_QUEUE_NAME"`
	EventLocalCreate                   bool     `env:"CONFIG_EVENTBUS_LOCAL_CREATE"`
	EventInMemory                      bool     `env:"CONFIG_EVENTBUS_INMEMORY"`
	EventEndpoint                      string   `env:"CONFIG_EVENTBUS_ENDPOINT"`
	EventSFNEndpoint                   string   `env:"CONFIG_EVENTBUS_SFN_ENDPOINT"`
	EventDynamicTaskARN                string   `env:"CONFIG_DYNAMIC_TASK_ARN"`
	EventTopics                        []string `env:"CONFIG_CONSUMES_TOPICS"`
	EventLargeMessageBucketLocalCreate bool     `env:"CONFIG_EVENTBUS_LARGE_MESSAGE_BUCKET_LOCAL_CREATE"`

	EventSQSQueueURL string `env:"CONFIG_EVENT_BUS_SQS_QUEUE_URL"`

	ObjstoreInMemory          bool     `env:"CONFIG_RESOURCES_OBJECT_STORE_INMEMORY"`
	ObjstoreLocalCreate       bool     `env:"CONFIG_RESOURCES_OBJECT_STORE_LOCAL_CREATE"`
	ObjstoreEndpoint          string   `env:"CONFIG_RESOURCES_OBJECT_STORE_ENDPOINT"`
	ObjstoreBuckets           []string `env:"CONFIG_RESOURCES_OBJECT_STORE_BUCKETS"`
	ObjstoreNoCopyQueryEscape bool     `env:"CONFIG_RESOURCES_OBJECT_NO_COPY_QUERY_ESCAPE"`

	ProjectionsEndpoint        string   `env:"CONFIG_RESOURCES_PROJECTIONS_ENDPOINT"`
	ProjectionsInMemory        bool     `env:"CONFIG_RESOURCES_PROJECTIONS_INMEMORY"`
	ProjectionsLocalCreate     bool     `env:"CONFIG_RESOURCES_PROJECTIONS_LOCAL_CREATE"`
	ProjectionsTables          []string `env:"CONFIG_RESOURCES_PROJECTIONS_TABLES"`
	ProjectionsCompositeTables []string `env:"CONFIG_RESOURCES_PROJECTIONS_COMPOSITE_TABLES"`

	ServiceHostnameFormat   string `env:"CONFIG_APP_SERVICE_HOSTNAME_FORMAT"`
	ServiceHostnameOverride string `env:"CONFIG_APP_SERVICE_HOSTNAME_OVERRIDE"`

	ShouldRunConsumer      bool `env:"CONFIG_APP_RUN_CONSUMER"`
	ShouldRunQueueConsumer bool `env:"CONFIG_APP_RUN_QUEUE_CONSUMER"`
	ShouldRunHTTP          bool `env:"CONFIG_APP_RUN_HTTP"`

	ShouldRunProducer   bool
	ShouldRunObjstore   bool
	ShouldRunHTTPClient bool

	EnvVars map[string]string
}

// ResolveConfigAndSetDefault will set the global config settings at runtime
// This should only be called once at start-up and then 'DefaultConfig' should be
// used to get the resolved config
func ResolveConfigAndSetDefault(opts ...Option) (Config, error) {
	var err error
	cfg, err = ResolveConfig(opts...)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// ResolveConfig will determine a runtime strategy and use that
// to resolve the runtime configuration
func ResolveConfig(opts ...Option) (Config, error) {
	o := &options{}
	o.apply(opts...)

	switch true {
	case o.FromENV:
		return resolveConfigFromENV(o)
	case o.FromServiceJSON:
		return resolveConfigFromServiceJSON(o)
	}
	return Config{}, fmt.Errorf("unable to determine runtime strategy")
}

// resolveConfigFromENV will create a runtime configuration using
// values set in ENV
func resolveConfigFromENV(opts *options) (Config, error) {
	mu.Lock()
	defer mu.Unlock()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to resolve env config: %w", err)
	}
	if opts != nil {
		cfg.ShouldRunProducer = opts.ShouldRunProducer
		cfg.ShouldRunObjstore = opts.ShouldRunObjstore
		cfg.ShouldRunHTTPClient = opts.ShouldRunHTTPClient
	} else {
		cfg.ShouldRunProducer = true
		cfg.ShouldRunObjstore = true
		cfg.ShouldRunHTTPClient = true
	}
	return cfg, nil
}

// resolveConfigFromServiceJSON will use the service.json to
// generate a runtime configuration
func resolveConfigFromServiceJSON(opts *options) (Config, error) {
	si, err := service.Read(opts.Fs, opts.ServiceJSONName)
	if err != nil {
		return Config{}, err
	}

	opts.ShouldRunProducer = si.EventBus.Producer
	opts.ShouldRunHTTPClient = si.HTTPClient

	cfg, err := resolveConfigFromENV(opts)
	if err != nil {
		return cfg, err
	}

	cfg.ShouldRunConsumer = si.EventBus.Consumer
	cfg.ShouldRunHTTP = si.HTTPS.Enabled

	cfg.ServiceName = si.Name
	cfg.HTTPHandlerTimeout = time.Duration(si.HTTPS.TimeoutSecond) * time.Second
	cfg.EventSQSQueueName = "event-bus-" + si.Name

	var topics []string
	for _, ct := range si.EventBus.ConsumesTopics {
		topics = append(topics, ct.Name)
	}
	cfg.EventTopics = topics

	var tables []string
	for _, d := range si.Resources.Databases {
		if !d.Deleted {
			tables = append(tables, d.Name)
		}
	}
	cfg.ProjectionsTables = tables

	var compositeTables []string
	for _, d := range si.Resources.CompositeDatabases {
		if !d.Deleted {
			compositeTables = append(compositeTables, d.Name)
		}
	}
	cfg.ProjectionsCompositeTables = compositeTables

	var buckets []string
	for _, os := range si.Resources.ObjectStores {
		if !os.Deleted {
			buckets = append(buckets, os.Name)
		}
	}
	cfg.ShouldRunObjstore = len(buckets) > 0
	cfg.ObjstoreBuckets = buckets

	cfg.EnvVars = si.EnvVars.DevelopmentEnvVars
	cfg.Team = si.Team
	cfg.Flow = si.CapabilityFlow()
	cfg.Tier = si.Tier
	cfg.Production = si.Production
	cfg.UseGha = si.UseGha
	cfg.IsGhe = si.IsGhe
	cfg.DisableSentry = si.DisableSentry
	return cfg, nil
}
