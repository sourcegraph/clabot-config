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

	liblog := log.Init(log.Resource{
		// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
		Name:       os.Getenv("GITHUB_WORKFLOW"),
		Namespace:  "clabot-config",
		InstanceID: os.Getenv("GITHUB_RUN_ID"),
	}, sinks...)
	return liblog.Sync
}
