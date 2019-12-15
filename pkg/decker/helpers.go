package decker

import "os"

// loop over list of strings and return true if list contains a given string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// makeDirectory makes a directory at a given path if it does not already exist
func makeDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
