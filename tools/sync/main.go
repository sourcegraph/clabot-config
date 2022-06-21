package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sourcegraph/clabot-config/internal/clabot"
	"github.com/sourcegraph/clabot-config/internal/responses"
)

var pages = flag.Int("pages", 10, "pages of responses to check")

// sync will check the CLA responses form for github handles that are not yet in .clabot,
// and add them to the .clabot file.
func main() {
	flag.Parse()

	ctx := context.Background()

	conf, err := clabot.ParseConfig()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	existingHandles := make(map[string]struct{}, len(conf.Contributors))
	for _, handle := range conf.Contributors {
		existingHandles[handle] = struct{}{}
	}

	resps, err := responses.ListResponses(ctx, *pages)
	if err != nil {
		// EOF indicates we hit maximum pages
		if err != io.EOF {
			println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Reached maximum pages %d\n", *pages)
	}
	fmt.Printf("Found %d responses\n", len(resps))
	for _, resp := range resps {
		// If not yet in config, add it
		if _, exists := existingHandles[resp.GitHubHandle]; !exists {
			fmt.Printf("Adding contributor %q (%s)\n", resp.GitHubHandle, resp.Name)
			conf.Contributors = append(conf.Contributors, resp.GitHubHandle)
			existingHandles[resp.GitHubHandle] = struct{}{} // for deduplication
		}
	}

	if err := conf.Save(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Configuration is up to date!")
}
