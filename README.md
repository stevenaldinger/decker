# What is decker

# Purpose

`Decker` is a penetration testing orchestration framework. It leverages [HashiCorp Configuration Language 2](https://github.com/hashicorp/hcl2) (the same config language as [Terraform](https://github.com/hashicorp/terraform)) to allow declarative `penetration testing as code`, so your tests can be versioned, shared, reused, and collaborated on with your team or the community.

Example of a `decker` config file:

```hcl
// variables are pulled from environment
//   ex: DECKER_TARGET_HOST
// they will be available throughout the config files as var.*
//   ex: ${var.target_host}
variable "target_host" {
  type = "string"
}

// resources refer to plugins
// resources need unique names so plugins can be used more than once
// they are declared with the form: 'resource "plugin_name" "unique_name" {}'
// their outputs will be available to others using the form unique_name.*
//   ex: nmap.443
resource "nmap" "nmap" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "sslscan" "sslscan" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}
```

Run a plugin for each item in a list:

```hcl
variable "target_host" {
  type = "string"
}
resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  dns_server = "8.8.4.4"
}
resource "metasploit" "metasploit" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/portscan/tcp"
  options = {
    RHOSTS = "${each.key}/32"
    INTERFACE = "eth0"
  }
}
```

## Why the name decker?

My friend [Courtney](https://github.com/courtneymiller2010) came to the rescue when I was struggling to come up with a name and found [decker](http://www.catb.org/esr/sf-words/glossary.html#decker) in a [SciFi word glossary](http://www.catb.org/esr/sf-words/glossary.html)... and it sounded cool.

> A future cracker; a software expert skilled at manipulating cyberspace, especially at circumventing security precautions.

# Running an example config with docker

Two volumes are mounted:

1. Directory named `decker-reports` where `decker` will output a file for each plugin executed. The file's name will be `{unique_resource_name}.report.txt`.
2. `examples` directory containing `decker` config files. Mounting this volume allows you to write configs locally using your favorite editor and still run them within the container.

One environment variable is passed in:

1. `DECKER_TARGET_HOST`

This is referenced in the config files as `{var.target_host}`. `Decker` will loop through all environment variables named `DECKER_*`, stripping away the prefix and setting the rest to lowercase.

```
docker run -it --rm \
  -v "$(pwd)/decker-reports/":/tmp/reports/ \
  -v "$(pwd)/examples/":/decker-config/ \
  -e DECKER_TARGET_HOST=example.com \
 stevenaldinger/decker:latest decker ./decker-config/example.hcl
```

# Contributing

Contributions are very welcome and appreciated. See [docs/contributions.md](docs/contributions.md) for guidelines.

# Development

Using docker for development is recommended for a smooth experience. This ensures all dependencies will be installed and ready to go.

Refer to `Directory Structure` below for an overview of the go code.

## Quick Start

1. (on host machine): `make docker_build`
2. (on host machine): `make docker_run` (will start docker container and open an interactive bash session)
3. (inside container): `dep ensure -v`
3. (inside container): `make build_all`
4. (inside container): `make run`

## Initialize git hooks

Run `make init` to add a `pre-commit` script that will run [linting](https://github.com/golang/lint) and tests on each commit.

# Plugin Development

`Decker` itself is just a framework that reads config files, determines dependencies in the config files, and runs plugins in an order that ensures plugins with dependencies on other plugins (output of one plugin being an input for another) run after the ones they depend on.

The real power of `decker` comes from plugins. Developing a plugin can be as simple or as complex as you want it to be, as long as the end result is a `.so` file containing the compiled plugin code and a `.hcl` file in the same directory declaring the inputs the plugin is expecting a user to configure.

The recommended way to get started with `decker` plugin development is by cloning the [decker-plugin](https://github.com/stevenaldinger/decker-plugin) repository and following the steps in its documentation. It should only take you a few minutes to get a "Hello World" `decker` plugin running.

## Installing plugins

Right now, plugins are expected to be in a directory relative to wherever the `decker` binary is.

Plugins are expected to be at `<decker binary>/internal/app/decker/plugins/<plugin name>/<plugin name>.so`.

There should be an `HCL` file next to the `.so` file at `<decker binary>/internal/app/decker/plugins/<plugin name>/<plugin name>.hcl` that defines its inputs. Currently, only `string` inputs are supported. Each input should have an `input` block that looks like this:

```
input "my_input" {
  type = "string"
  default = "some default value"
}
```

# Directory Structure

```
.
├── build
│   ├── ci/
│   └── package/
├── cmd
│   ├── decker
│   │   └── main.go
│   └── README.md
├── deployments/
├── docs/
├── examples
│   └── example.hcl
├── githooks
│   ├── pre-commit
├── Gopkg.toml
├── internal
│   ├── app
│   │   └── decker
│   │       └── plugins
│   │           ├── a2sv
│   │           │   ├── a2sv.hcl
│   │           │   ├── main.go
│   │           │   └── README.md
│   │           └── ...
│   │               ├── main.go
│   │               ├── README.md
│   │               └── xxx.hcl
│   ├── pkg
│   │   ├── dependencies/
│   │   ├── hcl/
│   │   ├── paths/
│   │   ├── plugins/
│   │   └── reports/
│   └── README.md
├── LICENSE
├── Makefile
├── README.md
└── scripts
    ├── build-plugins.sh
    └── README.md
```

- [cmd/decker/main.go](cmd/decker/main.go) is the driver. Its job is to parse a given [config file](examples/), load the appropriate plugins based on the file's `resource` blocks, and run the plugins with the specified inputs.
- [examples](examples/) has a couple example configurations to get you started with `decker`. If you use the given [docker image](https://hub.docker.com/r/stevenaldinger/decker/), all dependencies should be installed for both config files and things should run smoothly.
- [internal/pkg](internal/pkg) is where most of the actual code is. It contains all the packages imported by [main.go](cmd/decker/main.go).
  * [dependencies](internal/pkg/dependencies) is responsible for building the plugin dependency graph and returning a topologically sorted array that ensures plugins are run in a working order.
  * [hcl](internal/pkg/hcl) is responsible for parsing HCL files, including creating evaluation contexts that let blocks properly decode when they depend on other plugin blocks.
  * [paths](internal/pkg/paths) is responsible for returning file paths for the `decker` binary, config files, plugin config files, and generated reports.
  * [plugins](internal/pkg/plugins) is responsible for determining if plugins are enabled and running them.
  * [reports](internal/pkg/reports) is responsible for writing reports to the file system.
- [internal/app/decker/plugins](internal/app/decker/plugins) are modular pieces of code written as [Golang plugins](https://golang.org/pkg/plugin/), implementing a simple interface that allows them to be loaded and called at run-time with inputs specified in the plugin's config file (also in HCL). An example can be found at [internal/app/decker/plugins/nslookup/nslookup.hcl](internal/app/decker/plugins/nslookup/nslookup.hcl).
- [decker config files](examples) offer a declarative way to write penetration tests. The manifests are written in [HashiCorp Configuration Language 2](https://godoc.org/github.com/hashicorp/hcl2/hcl)) and describe the set of plugins to be used in the test as well as their inputs.
