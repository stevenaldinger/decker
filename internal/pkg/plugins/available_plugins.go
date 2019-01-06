package plugins

import (
	"io/ioutil"
	"log"

	"github.com/stevenaldinger/decker/internal/pkg/paths"
)

// AvailablePlugins searches the plugin directory and returns names of each
// plugin found.
func AvailablePlugins() []string {
	pluginDir := paths.GetPluginDirectory()
	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Fatal(err)
	}

	var plugins = []string{}

	for _, f := range files {
		if f.IsDir() {
			plugins = append(plugins, f.Name())
		}
	}

	return plugins
}
