package main

import (
	"context"
	"flag"
	"io"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/clabot-config/internal/clabot"
	"github.com/sourcegraph/clabot-config/internal/responses"
)

var pages = flag.Int("pages", 10, "pages of responses to check")

// sync will check the CLA responses form for github handles that are not yet in .clabot,
// and add them to the .clabot file.
func main() {
	flag.Parse()

	ctx := context.Background()

	sync := initLogging()
	defer sync()
	logger := log.Scoped("sync", "tool to sync contributors")

	// Read config and parse our contributors
	conf, err := clabot.ParseConfig()
	if err != nil {
		logger.Fatal("ParseConfig", log.Error(err))
	}
	existingHandles := make(map[string]struct{}, len(conf.Contributors))
	for _, handle := range conf.Contributors {
		existingHandles[handle] = struct{}{}
	}

	// List CLA form responses
	resps, err := responses.ListResponses(ctx, *pages)
	if err != nil {
		// EOF indicates we hit maximum pages
		if err != io.EOF {
			logger.Fatal("ListResponses", log.Error(err))
		}
		logger.Info("reached maximum pages", log.Int("pages", *pages))
	}
	logger.Info("listed responses", log.Int("responses", len(resps)))

	// For each response, if not yet in config, add it
	var added int
	for _, resp := range resps {
		if _, exists := existingHandles[resp.GitHubHandle]; !exists {
			logger.Info("adding contributor",
				log.String("gitHubHandle", resp.GitHubHandle),
				log.String("name", resp.Name))

			conf.Contributors = append(conf.Contributors, resp.GitHubHandle)
			existingHandles[resp.GitHubHandle] = struct{}{} // for deduplication
			added += 1
		}
	}

	logger = logger.With(log.Int("added", added))
	if added > 0 {
		// Write updated configuration back
		if err := conf.Save(); err != nil {
			logger.Fatal("conf.Save", log.Error(err))
		}
		logger.Info("configuration is up to date")
	} else {
		logger.Info("no updates to make")
	}
}
