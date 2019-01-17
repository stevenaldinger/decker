package paths

import (
	"os"
)

func getEnv(varName, defaultVal string) string {
	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}

	return defaultVal
}

// GetReportFilePath is given a resource name and returns the path for its
// report file. The directory can be set using environment variable
// DECKER_REPORTS_DIR or will default to "/tmp/reports"
func GetReportFilePath(resourceName, extension string) string {
	reportsDir := getEnv("DECKER_REPORTS_DIR", "/tmp/reports")

	filePath := reportsDir + "/" + resourceName + ".report." + extension

	return filePath
}
