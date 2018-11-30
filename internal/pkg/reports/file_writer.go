package reports

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// WriteStringToFile takes a file path and a content string and writes the
// content to a file with the given path
func WriteStringToFile(filePath, str string) {
	f, err := os.Create(filePath)
	check(err)

	defer f.Close()

	// written, err := f.WriteString("writes\n")
	f.WriteString(str)
	// fmt.Printf("wrote string:\n", written)

	// flush
	f.Sync()
}
