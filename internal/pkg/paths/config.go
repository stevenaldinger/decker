package paths

import (
	"fmt"
	"os"
)

// GetConfigFilePath gets HCL config file path from CLI args or print usage and exit
func GetConfigFilePath() string {
	hclConfigFile := ""
	if len(os.Args) == 2 {
		hclConfigFile = os.Args[1]
	} else {
		fmt.Println("Usage: \"decker ./examples/exploit-example.hcl\"")
		os.Exit(1)
	}

	return hclConfigFile
}
