package main

import (
	"bufio"
	"fmt"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"strconv"
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

// https://www.dotnetperls.com/duplicates-go
func removeDuplicates(elements []string) []string {
  encountered := map[string]bool{}

  // Create a map of all unique elements.
  for v:= range elements {
      encountered[elements[v]] = true
  }

  // Place all keys from the map into a slice.
  result := []string{}
  for key, _ := range encountered {
      result = append(result, key)
  }
  return result
}

type plugin string

// --- wants input: ---
// inputsMap{
//   "host": "example.com",
// }
//
// --- gives output: ---
// resultsMap{
//  "raw_output": "...",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	var (
		cmdOut []byte
		err    error
		cmdArgs []string
		outputByLine    []string
		services []string
		servicesMap map[string] []string
	)

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}
	servicesMap = map[string][]string{}

	targetHost := decoder.GetString((*inputsMap)["host"])
	scanType := decoder.GetString((*inputsMap)["type"])

	cmdName := "nmap"
	if scanType == "protocol_detection" {
		cmdArgs = []string{"-O", "-oG", "-", targetHost}
	} else {
		cmdArgs = []string{targetHost}
	}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running nmap command: ", err)
		return
	}

	output := string(cmdOut)

	if scanType == "protocol_detection" {
		// Host: 93.184.216.34 ()	Status: Up
		// Host: 93.184.216.34 ()	Ports: 80/open/tcp//http///, 443/open/tcp//https///, 1935/closed/tcp//rtmp///	Ignored State: filtered (997)	Seq Index: 252	IP ID Seq: Randomized
		// # Nmap done at Mon Mar 11 19:02:34 2019 -- 1 IP address (1 host up) scanned in 8.76 seconds

		// fmt.Println("Output:", output)

		if outputByLine, err = stringToLines(output); err != nil {
			fmt.Fprintln(os.Stderr, "Error occured while converting nmap protocol to array of strings: ", err)
			return
		}

		splitOutput := strings.Split(outputByLine[2], "Ports: ")
		portsOpen := splitOutput[1]
		portsOpenArray := strings.Split(portsOpen, ",")

		for _, portInfo := range portsOpenArray {
			// fmt.Println("Port:", portInfo)

			portInfoSplit := strings.Split(portInfo, "/")
			// fmt.Println("Len(portInfoSplit):", len(portInfoSplit))
			parsedPort := portInfoSplit[0]
			parsedState := portInfoSplit[1]
			parsedTransportProtocol := portInfoSplit[2]
			parsedServiceProtocol := portInfoSplit[4]
			// fmt.Println("Port:", parsedPort)
			// fmt.Println("State:", parsedState)
			// fmt.Println("Transport:", parsedTransportProtocol)
			// fmt.Println("Service:", parsedServiceProtocol)
			portInfoMap := map[string]cty.Value{
				"port": encoder.StringVal(parsedPort),
				"state": encoder.StringVal(parsedState),
				"transport": encoder.StringVal(parsedTransportProtocol),
				"service": encoder.StringVal(parsedServiceProtocol),
			}

			(*resultsMap)[parsedPort] = encoder.MapVal(portInfoMap)
			// (*resultsMap)[parsedPort] = encoder.StringVal(parsedServiceProtocol)

			if _, ok := servicesMap[parsedServiceProtocol]; !ok {
				servicesMap[parsedServiceProtocol] = []string{}
			}
			servicesMap[parsedServiceProtocol] = removeDuplicates(append(servicesMap[parsedServiceProtocol], parsedPort))
			services = append(services, parsedServiceProtocol)
		}

		services = removeDuplicates(services)

		// set everything that's not open to false so value is defined in HCL configs
		for i := 1; i <= 30000; i++ {
			if _, ok := (*resultsMap)[strconv.Itoa(i)]; !ok {
				portInfoMap := map[string]cty.Value{
					"port": encoder.StringVal(strconv.Itoa(i)),
					"state": encoder.StringVal("closed"),
					"transport": encoder.StringVal(""),
					"service": encoder.StringVal(""),
				}

				(*resultsMap)[strconv.Itoa(i)] = encoder.MapVal(portInfoMap)
				// (*resultsMap)[strconv.Itoa(i)] = encoder.StringVal("")
			}
		}

		var servicesCty = []cty.Value{}
		for _, service := range services {
			servicesCty = append(servicesCty, encoder.StringVal(service))

			var portsWithServiceCty = []cty.Value{}
			for _, port := range servicesMap[service] {
				portsWithServiceCty = append(portsWithServiceCty, encoder.StringVal(port))
			}
			(*resultsMap)[service] = encoder.ListVal(portsWithServiceCty)
		}

		(*resultsMap)["services"] = encoder.ListVal(servicesCty)
	}
	// } else {
	// 	portInfoMap := map[string]cty.Value{
	// 		"port": encoder.StringVal(strconv.Itoa(i)),
	// 		"state": encoder.StringVal("closed"),
	// 		"transport": encoder.StringVal(""),
	// 		"service": encoder.StringVal(""),
	// 	}
	// 	(*resultsMap)[strconv.Itoa(i)] = encoder.MapVal(portInfoMap)
	// 	(*resultsMap)["protocol"] = encoder.StringVal("")
	// }

	(*resultsMap)["raw_output"] = encoder.StringVal(output)
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
