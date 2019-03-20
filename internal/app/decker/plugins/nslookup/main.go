package main

import (
	"bufio"
	"fmt"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"strings"
)

func stringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return lines, err
}

type plugin string

// --- wants input: ---
// inputsMap{
//   "host": "example.com",
//   "dns_server": "8.8.4.4".
//   "plugin_enabled": "true",
// }
//
// --- gives output: ---
// resultsMap{
//  "dns_server": "8.8.4.4",
//  "dns_address": "8.8.4.4#53",
//  "host_name": "example.com",
//  "ip_address": "172.217.11.142",
//  "raw_output": "...",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	var (
		cmdOut          []byte
		err             error
		outputByLineDNS []string
		outputByLine    []string
	)

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	targetHost := decoder.GetString((*inputsMap)["host"])
	dnsServer := decoder.GetString((*inputsMap)["dns_server"])

	cmdName := "nslookup"
	cmdArgs := []string{targetHost, dnsServer}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running nslookup command: ", err)
		return
	}

	output := string(cmdOut)

	// WILL BREAK IF NO ANSWER
	splitOutput := strings.Split(output, "Non-authoritative answer:")
	dnsServerInfo, outputAfterDNSInfo := splitOutput[0], splitOutput[1]

	if outputByLineDNS, err = stringToLines(dnsServerInfo); err != nil {
		fmt.Fprintln(os.Stderr, "Error occured while converting nslookup output to array of strings: ", err)
		return
	}

	if outputByLine, err = stringToLines(outputAfterDNSInfo); err != nil {
		fmt.Fprintln(os.Stderr, "Error occured while converting nslookup output to array of strings: ", err)
		return
	}

	// parse out DNS server info
	for _, line := range outputByLineDNS {
		if strings.Contains(line, "Server:") {
			serverIP := strings.TrimSpace(strings.Split(line, ":")[1])
			(*resultsMap)["dns_server"] = encoder.StringVal(serverIP)
		}

		if strings.Contains(line, "Address:") {
			serverAddress := strings.TrimSpace(strings.Split(line, ":")[1])
			(*resultsMap)["dns_address"] = encoder.StringVal(serverAddress)
		}
	}

	var ipAddList = []string{}
	// (*resultsListMap)["ip_address"] = ipAddList
	// parse out host info
	for _, line := range outputByLine {
		if strings.Contains(line, "Name:") {
			hostName := strings.TrimSpace(strings.Split(line, ":")[1])
			(*resultsMap)["host_name"] = encoder.StringVal(hostName)
		}

		if strings.Contains(line, "Address:") {
			hostAddress := strings.TrimSpace(strings.Split(line, ":")[1])
			ipAddList = append(ipAddList, hostAddress)
			// (*resultsListMap)["ip_address"] = append((*resultsListMap)["ip_address"], hostAddress)
		}
	}

	var ipAddListCty = []cty.Value{}
	for _, ipAdd := range ipAddList {
		ipAddListCty = append(ipAddListCty, encoder.StringVal(ipAdd))
	}

	(*resultsMap)["ip_address"] = encoder.ListVal(ipAddListCty)
	(*resultsMap)["raw_output"] = encoder.StringVal(output)
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
