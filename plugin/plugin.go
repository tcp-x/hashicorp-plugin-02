package main

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	shared "github.com/tcp-x/cd-plug-util" // Import the shared package
)

type CdExecutorPlugin struct {
	logger hclog.Logger
}

// CdExec is the method that executes a command and returns the result.
func (exec *CdExecutorPlugin) CdExec(jsonInput string) (string, error) {
	exec.logger.Debug("message from CdExecutorPlugin.CdExec")
	// Parse the JSON input
	var input map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &input)
	if err != nil {
		return "", err
	}

	// Your logic here to execute the command based on the input
	// For demonstration, let's just return a JSON response
	response := map[string]interface{}{
		"data":   "Your result data",
		"status": "success",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(responseJSON), nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	exec := &CdExecutorPlugin{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"cd_executor": &shared.CdExecutorPlugin{Impl: exec},
	}

	logger.Debug("message from plugin", "foo", "bar")
	// Serve the plugin using the CdExecutorPlugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
