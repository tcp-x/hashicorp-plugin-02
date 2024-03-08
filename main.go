package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	shared "github.com/tcp-x/cd-plug-util" // Import the shared package
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"cd_executor": &shared.CdExecutorPlugin{},
}

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// Create a plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/executor"),
		Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	// defer rpcClient.Close()
	defer client.Kill()

	// Request the plugin
	raw, err := rpcClient.Dispense("cd_executor")
	if err != nil {
		log.Fatal("Error getting plugin:", err)
	}
	executor := raw.(shared.CdExecutor)

	// Example input JSON string
	input := `{"key": "value"}`

	// Call the CdExec method of the plugin
	result, err := executor.CdExec(input)
	if err != nil {
		log.Fatal("Error calling CdExec:", err)
	}

	fmt.Println("Plugin Result:", result)
}
