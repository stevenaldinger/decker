package paths

// GetPluginHCLFilePath is given a plugin name and returns the path to its HCL
// config file (which defines its inputs and outputs).
func GetPluginHCLFilePath(pluginName string) string {
	// plugins are expected to be in a location relative to the decker binary
	deckerDir := GetDeckerDir()

	filePath := deckerDir + "/internal/app/decker/plugins/" + pluginName + "/" + pluginName + ".hcl"

	return filePath
}
