package main

import (
	"context"
	"strings"

	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"

	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
)

// version is set via ldflags by GoReleaser.
var version = "dev"

const pluginName = "echo"

// Config holds executor configuration.
type Config struct {
	ChangeResponseToUpperCase *bool `yaml:"changeResponseToUpperCase,omitempty"`
}

// EchoExecutor implements Botkube executor plugin.
type EchoExecutor struct{}

// Metadata returns details about Echo plugin.
func (EchoExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: "Echo is an example Botkube executor plugin used during e2e tests. It's not meant for production usage.",
	}, nil
}

// Execute returns a given command as response.
func (EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	// In our case we don't have complex merge strategy,
	// the last one that was specified wins :)
	finalCfg := Config{}
	for _, inputCfg := range in.Configs {
		var cfg Config
		err := yaml.Unmarshal(inputCfg.RawYAML, &cfg)
		if err != nil {
			return executor.ExecuteOutput{}, err
		}
		if cfg.ChangeResponseToUpperCase == nil {
			continue
		}
		finalCfg.ChangeResponseToUpperCase = cfg.ChangeResponseToUpperCase
	}

	data := in.Command
	if finalCfg.ChangeResponseToUpperCase != nil && *finalCfg.ChangeResponseToUpperCase {
		data = strings.ToUpper(data)
	}

	return executor.ExecuteOutput{
		Data: data,
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		pluginName: &executor.Plugin{
			Executor: &EchoExecutor{},
		},
	})
}
