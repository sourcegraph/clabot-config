package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/log"
)

func initLogging() func() error {
	var sinks []log.Sink
	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		sinks = append(sinks, log.NewSentrySinkWith(log.SentrySink{
			ClientOptions: sentry.ClientOptions{
				Dsn:        sentryDSN,
				SampleRate: 1, // send all
			},
		}))
	}

	liblog := log.Init(log.Resource{Name: "clabot-config"}, sinks...)
	return liblog.Sync
}
