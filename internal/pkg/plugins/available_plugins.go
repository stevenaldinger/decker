package plugins

import (
	"github.com/stevenaldinger/decker/internal/pkg/paths"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func getEnv(varName, defaultVal string) string {
	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}

	return defaultVal
}

// AvailablePlugins searches the plugin directories and returns names of each
// plugin found.
func AvailablePlugins() []string {
	var pluginDirs = []string{}
	var plugins = []string{}

	pluginDirs = append(pluginDirs, paths.GetPluginDirectory())

	dirs := strings.Split(getEnv("DECKER_PLUGIN_DIRS", ""), ":")

	for _, dir := range dirs {
		if len(dir) > 0 {
			pluginDirs = append(pluginDirs, dir)
		}
	}

	for _, pluginDir := range pluginDirs {
		files, err := ioutil.ReadDir(pluginDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.IsDir() {
				plugins = append(plugins, f.Name())
			}
		}
	}

	return plugins
}
