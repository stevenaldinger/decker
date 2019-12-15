package paths

import (
	"fmt"
	"os"
)

// GetConfigFilePath gets HCL config file path from CLI args or print usage and exit
func GetConfigFilePath() string {
	hclConfigFile := ""
	jsonHCL := os.Getenv("DECKER_RUN_CONFIGURATION")
	if len(os.Args) == 2 {
		hclConfigFile = os.Args[1]
	} else if len(jsonHCL) == 0 {
		fmt.Println("Usage: \"decker ./examples/exploit-example.hcl\"")
		os.Exit(1)
	}

	return hclConfigFile
}
