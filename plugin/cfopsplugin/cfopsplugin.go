package cfopsplugin

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/xchapter7x/lo"
)

//Start - takes a given plugin and starts it
func Start(plgn Plugin) {
	gob.Register(plgn)
	RegisterPlugin(plgn.GetMeta().Name, plgn)

	if len(os.Args) == 2 && os.Args[1] == PluginMeta {
		b, _ := json.Marshal(plgn.GetMeta())
		UIOutput(string(b))

	} else {

		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         GetPlugins(),
		})
	}
}

//Call - calls a given plugin by name and file path
func Call(name string, filePath string) (BackupRestorer, *plugin.Client) {
	var backupRestoreInterface BackupRestorer
	RegisterPlugin(name, backupRestoreInterface)
	log.SetOutput(ioutil.Discard)

	client := plugin.NewClient(&plugin.ClientConfig{
		Stderr:          os.Stderr,
		SyncStdout:      os.Stdout,
		SyncStderr:      os.Stderr,
		HandshakeConfig: GetHandshake(),
		Plugins:         GetPlugins(),
		Cmd:             exec.Command(filePath, "plugin"),
	})

	rpcClient, err := client.Client()
	if err != nil {
		lo.G.Debug("rpcclient error: ", err)
		log.Fatal(err)
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		lo.G.Debug("error: ", err)
		log.Fatal(err)
	}
	return raw.(BackupRestorer), client
}

//RegisterPlugin - register a plugin as available
func RegisterPlugin(name string, plugin BackupRestorer) {
	lo.G.Debug("registering plugin: ", name, plugin)
	pluginMap[name] = &BackupRestorePlugin{
		P: plugin,
	}
}

//GetPlugins - returns the list of registered plugins
func GetPlugins() map[string]plugin.Plugin {
	return pluginMap
}

//GetHandshake - gets the handshake object
func GetHandshake() plugin.HandshakeConfig {
	return handshakeConfig
}
