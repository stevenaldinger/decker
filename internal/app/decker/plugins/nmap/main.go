package main

import (
	"github.com/t94j0/nmap"
	"strconv"
	"strings"
)

type plugin string

// --- wants input: ---
// inputsMap{
//   "host": "example.com",
// }
//
// --- gives output: ---
// resultsMap{
//  "host_address": "123.456.789.1",
//  "host": "example.com",
//  "22": "open",
//  "80": "open",
//  "443": "open",
//  "raw_output": "...",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]string, resultsListMap *map[string][]string) {
	targetHost := strings.TrimSpace((*inputsMap)["host"])

	scan, _ := nmap.Init().AddHosts(targetHost).Run()
	host, _ := scan.GetHost(targetHost)

	for _, port := range host.Ports {
		(*resultsMap)[strconv.Itoa(int(port.ID))] = port.State
	}

	(*resultsMap)["host_address"] = host.Address

	if len(host.Hostnames) > 0 {
		(*resultsMap)["host"] = host.Hostnames[0].Name
	}
	(*resultsMap)["raw_output"] = host.ToString()
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
