# Building Decker Plugins

A "hello world" example plugin can be found at [internal/app/decker/plugins/hello_world](../internal/app/decker/plugins/hello_world).

## Run a decker kali docker container for development

```sh
make docker_run
```

This will start up a decker docker container with the decker repo's code mounted so you can make changes in your editor and compile and run them within the container.

## Compile decker and all plugins

Once you're inside the container, run this to compile decker and all of the plugins:

```sh
make build_all
```

## Compile a single plugin

This will build the plugin found at [internal/app/decker/plugins/hello_world](../internal/app/decker/plugins/hello_world).

```sh
make plugin=hello_world build_plugin
```

## Plugin Guide

Plugins are passed in an inputs map and outputs map (both `*map[string]cty.Value`) that `decker` creates automatically using inputs from the plugin's `resource block` in an `hcl` config file. In [examples/hello-world.hcl](../examples/hello-world.hcl), the resource block is:

```hcl
resource "hello_world" "my_plugin" {
  say_hello_to = "${var.say_hello}"
}
```

This will run the plugin at [internal/app/decker/plugins/hello_world](../internal/app/decker/plugins/hello_world) and the plugin's output will be written to `/tmp/reports/my_plugin.report.txt`.

Decker knows to look for `say_hello_to` in the resource block because an input is defined in [internal/app/decker/plugins/hello_world/hello_world.hcl](../internal/app/decker/plugins/hello_world/hello_world.hcl).

```hcl
input "say_hello_to" {
  type = "string"
  default = "world"
}
```

To try to keep the plugin's Golang interface simple while supporting more complex inputs, [go-cty](https://github.com/zclconf/go-cty) values are used which allows all types to be passed into a single input/output map. Decker offers helpers for encoding/decoding go-cty and native Golang types.

This will be passed into the plugin so it can be used in the plugin's code like this:

```go
sayHelloTo := decoder.GetString((*inputsMap)["say_hello_to"])
```

Outputs do not currently need to be defined in the plugin's hcl the way that inputs do. Anything that the plugin's code adds to the `resultsMap` will be available to use in the decker config files.

For example, we can take the value that was passed in as an input and output it as `said_hello_to`. In the code, that looks like this:

```go
(*resultsMap)["said_hello_to"] = encoder.StringVal(sayHelloTo)
```

This can be passed into another plugin in a decker config file like this:

```hcl
resource "hello_world" "my_plugin_2" {
  say_hello_to = "${my_plugin.said_hello_to}... again!"
}
```

This plugin's output can be found at `/tmp/reports/my_plugin_2.report.txt`. If `world` was passed into `my_plugin`, the output of `my_plugin_2` would be `Hello, world... again!`.

Refer to the [internal/app/decker/plugins](../internal/app/decker/plugins) for examples of using the helpers.

An example of outputting lists can be found in [internal/app/decker/plugins/nslookup](../internal/app/decker/plugins/nslookup):

```go
var ipAddListCty = []cty.Value{}
for _, ipAdd := range ipAddList {
  ipAddListCty = append(ipAddListCty, encoder.StringVal(ipAdd))
}

(*resultsMap)["ip_address"] = encoder.ListVal(ipAddListCty)
```

An example using a map/dictionary input can be found in [internal/app/decker/plugins/metasploit](../internal/app/decker/plugins/metasploit):

```go
	options := decoder.GetMap((*inputsMap)["options"])
  for key, val := range *options {
    cmdStr = cmdStr + "set " + key + " " + val + ";"
  }
```

The input declaration for a map does not need to specify the keys that will be inside of it, the type just needs to be map.

```hcl
input "options" {
  type = "map"
  default = {}
}
```

## Running a config file

The [examples/hello-world.hcl](../examples/hello-world.hcl) decker config file can be run either with:

```sh
make run_hello_world
```

or with:

```sh
./decker ./examples/hello-world.hcl
```
