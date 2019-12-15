package main

import (
	"github.com/stevenaldinger/decker/pkg/decker"
)

func main() {
	app := decker.Init()

	// run plugins
	app.RunPlugins()

	// run outputs plugin
	app.GetResults()
}
