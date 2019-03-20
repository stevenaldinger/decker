package main

import (
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/t94j0/nmap"
	"github.com/zclconf/go-cty/cty"
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
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	targetHost := strings.TrimSpace(decoder.GetString((*inputsMap)["host"]))

	scan, _ := nmap.Init().AddHosts(targetHost).Run()
	host, _ := scan.GetHost(targetHost)

	for _, port := range host.Ports {
		(*resultsMap)[strconv.Itoa(int(port.ID))] = encoder.StringVal(port.State)
	}

	// set everything that's not open to false so value is defined in HCL configs
	for i := 1; i <= 30000; i++ {
		if _, ok := (*resultsMap)[strconv.Itoa(i)]; !ok {
			(*resultsMap)[strconv.Itoa(i)] = encoder.StringVal("closed")
		}
	}

	(*resultsMap)["host_address"] = encoder.StringVal(host.Address)

	if len(host.Hostnames) > 0 {
		(*resultsMap)["host"] = encoder.StringVal(host.Hostnames[0].Name)
	}
	(*resultsMap)["raw_output"] = encoder.StringVal(host.ToString())
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
