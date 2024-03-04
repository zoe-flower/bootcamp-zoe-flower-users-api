// Code generated by dev-tool; DO NOT EDIT.
package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/flypay/events/pkg/connector"

	"github.com/flypay/go-kit/v4/pkg/app"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/flypay/go-kit/v4/pkg/safe"
	"github.com/flypay/go-kit/v4/pkg/tracing"

	service "github.com/flypay/bootcamp-zoe_flower-users-api/internal/app"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/getsentry/sentry-go"

	_ "go.uber.org/automaxprocs"
)

// NOTE: not used anymore, everything comes through ENVs
var runOnlyFlag string

func init() {
	flag.StringVar(&runOnlyFlag, "run-only", "", "Run only a specific API within this service")
}

func main() {
	flag.Parse()

	opts := []runtime.Option{
		runtime.WithFromENV(),
		runtime.WithShouldRunProducer(ConnectorDefinition.Identifier != ""),
	}

	cfg, err := runtime.ResolveConfigAndSetDefault(opts...)
	if err != nil {
		log.Fatal("unable to resolve runtime config: %v", err)
	}

	log.SetLogger(
		os.Stdout,
		cfg.ServiceName,
		cfg.Version,
		cfg.Team,
		cfg.Flow,
		cfg.Tier,
		cfg.Production,
	)
	log.Info("initialised logger")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDSN,
		Release:     cfg.Version,
		Environment: cfg.Environment,
	}); err != nil {
		log.Fatalf("Sentry init failed: %+v", err)
	}
	defer sentry.Flush(time.Second * 5)
	log.SetSentryHook(sentry.CurrentHub().Client())

	defer safe.Recover()

	log.Info("Initialised sentry")

	if err := xray.Configure(xray.Config{
		DaemonAddr: cfg.XrayDaemon,
	}); err != nil {
		log.Fatalf("XRAY init failed: %+v", err)
	}

	if cfg.XrayLocalLogger {
		// This is to redirect xray to the Flyt logger when running locally
		xray.SetLogger(&log.XrayLocalLogger{Logger: log.DefaultLogger})
	}

	privhttpsrv := app.StartInternalHTTP()

	log.Info("Started http server for liveness, readiness and metrics")

	runner, err := app.NewRunner(cfg)
	if err != nil {
		log.Fatalf("Unable to create new service runner: %v", err)
	}

	if err := service.RunService(

		runner.Echo,
	); err != nil {
		log.Fatalf("Failed to init service: %+v", err)
	}

	safe.Go(func() {
		if ConnectorDefinition.Identifier != "" {
			produceEvent(runner.Producer, cfg.ServiceName)
		}
	})

	if cfg.ShouldRunHTTP {
		safe.Go(func() {
			log.Infof("Starting external http server on http://localhost:%s", cfg.HTTPPort)
			log.Error(runner.Echo.Start(":" + cfg.HTTPPort))
		})
	}

	app.GracefullyShutdown(privhttpsrv)
	app.Run()
}

func produceEvent(producer eventbus.Producer, appName string) {
	for {
		ctx, seg := tracing.StartSpan(context.Background(), appName)

		err := producer.Emit(ctx, &connector.Published{Connector: &ConnectorDefinition})
		if err != nil {
			log.Errorf("Failed to emit event: %+v", err)
			seg.Close(err)
			time.Sleep(2 * time.Second)
			continue
		}

		seg.Close(nil)
		break
	}

	log.Info("Emitted connector.Published event")
}
