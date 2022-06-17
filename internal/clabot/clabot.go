package clabot

import (
	"encoding/json"
	"os"
	"sort"
)

const configPath = ".clabot"

type Config struct {
	Message        string   `json:"message"`
	Label          string   `json:"label"`
	RecheckComment string   `json:"recheckComment"`
	Contributors   []string `json:"contributors"`
}

func ParseConfig() (*Config, error) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	return &config, json.Unmarshal(b, &config)
}

func (c *Config) Save() error {
	// Ensure consistent output
	sort.Strings(c.Contributors)

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, append(b, '\n'), os.ModePerm)
}
