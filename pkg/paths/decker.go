package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDeckerDir finds the directory the "decker" binary is in.
func GetDeckerDir() string {
	deckerPath, err := os.Executable()

	if err != nil {
		fmt.Println("Error finding executable path:", err)
		os.Exit(1)
	}

	deckerDir := filepath.Dir(deckerPath)

	return deckerDir
}
