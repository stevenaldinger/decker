# Building Go Plugins

```sh
CGO_ENABLED=1 go build -buildmode=plugin -o plugin.so "$plugins_dir/main.go"`
```

# Inputs

Plugins are passed in an inputs map (`*map[string]string`) that `decker` makes using inputs from the plugin's `resource block` in an `hcl` config file.

To try to keep the interface simple while supporting more complex inputs, arrays and maps are converted to JSON strings so they can be added to the inputs map.

## Example supporting array input

Top level config file:

```hcl
resource "my_plugin" "my_plugin" {
  example_input_array = ["input 1", "input 2"]
}
```

Plugin config file:

```hcl
input "example_input_array" {
  type = "list"
  default = []
}
```

Plugin code:

```go
package main

import (
  ...
  "encoding/json"
)

func (p plugin) Run(inputsMap, resultsMap *map[string]string) {
  var exampleInputs []string

  err := json.Unmarshal([]byte((*inputsMap)["example_input_array"]), &exampleInputs)

  ...
}
```

## Example supporting map input

Top level config file:

```hcl
resource "my_plugin" "my_plugin" {
  example_input_map = {
    input_1 = "some_input"
    input_2 = "some_other_input"
  }
}
```

Plugin config file:

```hcl
input "example_input_map" {
  type = "map"
  default = {}
}
```

Plugin code:

```go
package main

import (
  ...
	"encoding/json"
)

func (p plugin) Run(inputsMap, resultsMap *map[string]string) {
	var exampleInput map[string]interface{}

	err := json.Unmarshal([]byte((*inputsMap)["example_input_map"]), &exampleInput)

  ...
}
```
