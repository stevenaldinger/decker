package paths

import (
	"io/ioutil"
	"log"
	"strings"
)

// GetPluginPath takes the name of a plugin, searches the plugin directories
// for a match, and returns the path to the plugin or empty string if the
// plugin isn't found
func GetPluginPath(pluginName string) string {
	var pluginDirs = []string{}

	pluginDirs = append(pluginDirs, GetPluginDirectory())

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
			if f.IsDir() && f.Name() == pluginName {
				pluginPath := pluginDir + "/" + f.Name()
				return pluginPath
			}
		}
	}

	return ""
}

// GetPluginHCLFilePath is given a plugin name and returns the path to its HCL
// config file (which defines its inputs and outputs).
func GetPluginHCLFilePath(pluginName string) string {
	filePath := GetPluginPath(pluginName) + "/" + pluginName + ".hcl"

	return filePath
}

// GetPluginDirectory gets the default directory decker searches for plugins.
func GetPluginDirectory() string {
	deckerDir := GetDeckerDir()

	pluginDir := deckerDir + "/internal/app/decker/plugins"

	return pluginDir
}
