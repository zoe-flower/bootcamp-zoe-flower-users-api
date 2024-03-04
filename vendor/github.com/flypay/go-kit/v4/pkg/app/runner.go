package app

import (
	"fmt"
	"net/http"

	"github.com/flypay/go-kit/v4/internal/service"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/flypay/events/pkg/flyt"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/amazon"
	"github.com/flypay/go-kit/v4/pkg/eventbus/amazon/custom"
	"github.com/flypay/go-kit/v4/pkg/eventbus/amazon/custom/customfakes"
	"github.com/flypay/go-kit/v4/pkg/eventbus/eventbusfakes"
	"github.com/flypay/go-kit/v4/pkg/eventbus/middleware"
	"github.com/flypay/go-kit/v4/pkg/flythttp"
	httpmiddleware "github.com/flypay/go-kit/v4/pkg/flythttp/server/middleware"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/objstore"
	"github.com/flypay/go-kit/v4/pkg/projections"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"google.golang.org/protobuf/encoding/protojson"
)

// Runner contains the resources required to run a s service. The resources
// from this can be fed directly into a services 'RunService'
type Runner struct {
	cfg  runtime.Config
	Name string

	s3Session *session.Session
	S3Store   objstore.DownloadUploader
	Echo      *echo.Echo

	Producer      eventbus.Producer
	Consumer      eventbus.Consumer
	queueConsumer custom.QueueConsumer

	HTTPClient *http.Client

	EnvVars map[string]string
}

// NewRunner will configure the resources to be used in a services 'RunService'.
// This should not be called from a service, but instead be called when wanting setup
// and run a service, ie: a services main.gen.go file or monorepo
func NewRunner(cfg runtime.Config) (*Runner, error) {
	tracing := middleware.NewTracing(cfg.ServiceName, cfg.Team, cfg.Flow, cfg.Tier)
	testing := middleware.NewTesting()
	logging := middleware.NewLogging()
	requestIDs := middleware.NewRequestTracing()
	flytLocationID := middleware.NewFlytLocationIDMigration()
	metrics := middleware.NewMetrics(middleware.Clock{}, flyt.E_TopicName, cfg.EventWorkerPool)
	sentryMiddleware := middleware.NewSentry(sentry.CurrentHub())
	executionID := middleware.NewExecutionTracing()

	snssqsSession, err := amazon.MakeEnvAwareSQSSNSSessionWithRuntimeConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create an sqs/sns session: %w", err)
	}

	var consumer eventbus.Consumer
	if cfg.ShouldRunConsumer {
		log.Info("starting the SQS event bus consumer")

		s3Session, err := objstore.MakeEnvAwareSessionWithRuntimeConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create an s3 session: %w", err)
		}

		queueURL := amazon.BuildQueueURL(cfg.AWSRegion, cfg.AWSAccountID, cfg.EventSQSQueueName)
		consumer, err = amazon.ResolveNewConsumerWithRuntimeConfig(amazon.ConsumerConfig{
			Sessions: amazon.ConsumerSessions{
				SQSSNSSession: snssqsSession,
				S3Session:     s3Session,
			},
			Region:                     cfg.AWSRegion,
			AccountID:                  cfg.AWSAccountID,
			ServiceName:                cfg.ServiceName,
			Queue:                      cfg.EventSQSTopicName,
			QueueURL:                   queueURL,
			DefaultEventHandlerTimeout: cfg.EventHandlerTimeout,
			WorkerPoolSize:             cfg.EventWorkerPool,
			VisibilityTimeout:          cfg.EventHandlerVisibilityTimeout,
			TopicName:                  eventbus.ProtoTopicNamingStrategy(flyt.E_TopicName),
			TopicARN:                   amazon.FlytTopicARNStrategy,
		}, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create SQS consumer: %w", err)
		}
		// Ensure the context has values from headers first
		consumer.Use(requestIDs.ConsumerMiddleware)
		consumer.Use(executionID.ConsumerMiddleware)

		// Now we can log with valid context
		consumer.Use(logging.ConsumerMiddleware)
		consumer.Use(metrics.ConsumerMiddleware)

		// Order not important
		consumer.Use(tracing.ConsumerMiddleware)
		consumer.Use(testing.ConsumerMiddleware)
		consumer.Use(flytLocationID.ConsumerMiddleware)

		// Must be last to ensure application panics are converted to errors
		consumer.Use(sentryMiddleware.ConsumerMiddleware)
	} else {
		consumer = &eventbusfakes.FakeConsumer{}
	}

	// Create the queue consumer if we should, and then return the instantiated one
	// later
	var queueConsumer custom.QueueConsumer
	if cfg.ShouldRunQueueConsumer {
		amazonQueueConsumer, err := custom.NewConsumer(custom.ConsumerConfig{
			Session:                    snssqsSession,
			ServiceName:                cfg.ServiceName,
			Queue:                      cfg.EventSQSTopicName,
			QueueURL:                   cfg.EventSQSQueueURL,
			DefaultEventHandlerTimeout: cfg.EventHandlerTimeout,
			WorkerPoolSize:             cfg.EventWorkerPool,
			VisibilityTimeout:          cfg.EventHandlerVisibilityTimeout,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create SQS/SNS queue consumer: %w", err)
		}

		metrics := middleware.NewMetricsForQueue(middleware.Clock{}, cfg.EventSQSTopicName, cfg.EventWorkerPool)
		amazonQueueConsumer.Use(metrics.ConsumerMiddleware)
		amazonQueueConsumer.Use(tracing.ConsumerMiddleware)
		amazonQueueConsumer.Use(sentryMiddleware.ConsumerMiddleware)
		queueConsumer = amazonQueueConsumer
	} else {
		queueConsumer = &customfakes.FakeQueueConsumer{}
	}

	var producer eventbus.Producer
	if cfg.ShouldRunProducer {
		log.Info("initialising SNS producer")

		s3Session, err := objstore.MakeEnvAwareSessionWithRuntimeConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create an s3 session session: %w", err)
		}

		sfnSession, err := amazon.MakeEnvAwareSFNSessionWithRuntimeConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create an sfn session: %w", err)
		}

		producer, err = amazon.ResolveNewProducerWithRuntimeConfig(amazon.ProducerConfig{
			Sessions: amazon.ProducerSessions{
				SQSSNSSession: snssqsSession,
				SFNSession:    sfnSession,
				S3Session:     s3Session,
			},
			AccountID:            cfg.AWSAccountID,
			Region:               cfg.AWSRegion,
			DelayStateMachineARN: cfg.EventDynamicTaskARN,
			TopicName:            eventbus.ProtoTopicNamingStrategy(flyt.E_TopicName),
			TopicARN:             amazon.FlytTopicARNStrategy,
			MarshallerOptions:    func(opt *protojson.MarshalOptions) { opt.UseEnumNumbers = true },
			LargeMessageBucket:   fmt.Sprintf("%s-%s-%s", cfg.Environment, cfg.ServiceName, service.LargeMessagesStore(cfg.ServiceName).Bucket),
		}, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create SNS producer: %w", err)
		}

		producer.Use(tracing.ProducerMiddleware)
		producer.Use(testing.ProducerMiddleware)
		producer.Use(metrics.ProducerMiddleware)
		producer.Use(logging.ProducerMiddleware)
		producer.Use(requestIDs.ProducerMiddleware)
		producer.Use(flytLocationID.ProducerMiddleware)
	} else {
		producer = &eventbusfakes.FakeProducer{}
	}

	var (
		storeSession *session.Session
		store        objstore.DownloadUploader
	)
	if cfg.ShouldRunObjstore {
		log.Info("connecting to S3 object storage")
		storeSession, err = objstore.MakeEnvAwareSessionWithRuntimeConfig(cfg)
		if err != nil {
			log.Fatalf("failed to create object storage session: %+v", err)
		}

		store, err = objstore.ResolveObjectStorageWithRuntimeConfig(storeSession, cfg)
		if err != nil {
			log.Fatalf("failed to create object storage: %+v", err)
		}
	}

	log.Info("initialising external http server")
	pubhttpsrv := echo.New()
	pubhttpsrv.HideBanner = true

	if cfg.ShouldRunHTTP {
		log.Info("applying middleware to http server")

		// Pre Middleware (Runs before Router)
		pubhttpsrv.Pre(httpmiddleware.XRay(cfg.ServiceName))

		// Middleware run after Router
		// The logger middleware swallows any returned errors, so this needs to be
		// the last called in the middleware chain.
		pubhttpsrv.Use(echomiddleware.Logger())
		pubhttpsrv.Use(httpmiddleware.Metrics())
		pubhttpsrv.Use(echomiddleware.TimeoutWithConfig(echomiddleware.TimeoutConfig{
			Timeout: cfg.HTTPHandlerTimeout,
		}))
		pubhttpsrv.Use(httpmiddleware.AddRequestID())
		pubhttpsrv.Use(httpmiddleware.AddJobID())
		pubhttpsrv.Use(httpmiddleware.AddExecutionID())
		pubhttpsrv.Use(httpmiddleware.ExecutionTracing())
		pubhttpsrv.Use(httpmiddleware.CacheControl)

		// needs to be the last in the middleware chain
		pubhttpsrv.Use(httpmiddleware.Recover(sentry.CurrentHub(), cfg.PanicOnRecover))

		// Default endpoint
		pubhttpsrv.GET("/ping", func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "Pong")
		})
	}

	// Add the resources to the global graceful shutdown
	eventbusGracefulClose := eventbus.GracefulClose{
		Consumers: []eventbus.Consumer{consumer, queueConsumer},
		Producers: []eventbus.Producer{producer},
	}
	GracefullyShutdown(pubhttpsrv, eventbusGracefulClose)

	return &Runner{
		cfg:           cfg,
		Name:          cfg.ServiceName,
		S3Store:       store,
		s3Session:     storeSession,
		Echo:          pubhttpsrv,
		Producer:      producer,
		Consumer:      consumer,
		queueConsumer: queueConsumer,
		HTTPClient:    flythttp.NewClientWithXray(),
		EnvVars:       cfg.EnvVars,
	}, nil
}

func (r *Runner) NewDynamoStorage(name string) projections.ReadWriter {
	log.Infof("connecting to dynamo database %s", name)
	sess, err := projections.MakeEnvAwareSessionWithRuntimeConfig("", r.cfg)
	if err != nil {
		log.Fatalf("failed to create a dynamoDB session: [%s] %s", name, err)
	}
	db, err := projections.ResolveDynamoDbWithRuntimeConfig(fmt.Sprintf("%s-%s", r.Name, name), sess, r.cfg)
	if err != nil {
		log.Fatalf("failed to create the dynamoDB client: [%s] %s", name, err)
	}
	return db
}

func (r *Runner) NewDynamoCompositeStorage(name string) projections.CompositeReadWriter {
	log.Infof("connecting to composite dynamo database %s", name)
	sess, err := projections.MakeEnvAwareSessionWithRuntimeConfig("", r.cfg)
	if err != nil {
		log.Fatalf("failed to create a dynamoDB session: [%s] %s", name, err)
	}
	db, err := projections.ResolveDynamoCompositeDbWithRuntimeConfig(fmt.Sprintf("%s-composite-%s", r.Name, name), sess, r.cfg)
	if err != nil {
		log.Fatalf("failed to create the dynamoDB client: [%s] %s", name, err)
	}
	return db
}

// CreateBuckets will create all buckets locally
func (r *Runner) CreateBuckets(environment, service string, buckets []string) error {
	return objstore.CreateLocalS3Buckets(r.s3Session, environment, service, buckets)
}

func (r *Runner) NewQueueConsumer(name string) custom.QueueConsumer {
	// Only return the amazon consumer when this pod is the correct one
	if r.cfg.ShouldRunQueueConsumer && r.cfg.EventSQSTopicName == name {
		log.Infof("starting the queue SQS/SNS event bus consumer for %s", name)
		return r.queueConsumer
	}
	log.Infof("%s using fake queue consumer", name)
	return &customfakes.FakeQueueConsumer{}
}
