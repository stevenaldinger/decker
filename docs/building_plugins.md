# Building Go Plugins

```sh
CGO_ENABLED=1 go build -buildmode=plugin -o plugin.so "$plugins_dir/main.go"`
```

# Inputs

Plugins are passed in an inputs map and outputs map (both `*map[string]cty.Value`) that `decker` makes using inputs from the plugin's `resource block` in an `hcl` config file.

To try to keep the interface simple while supporting more complex inputs, [go-cty](https://github.com/zclconf/go-cty) values are used which allows all types to be passed into a single input/output map. Decker offers helpers for encoding/decoding go-cty and native Golang types. Refer to the [existing plugins](../internal/app/decker/plugins) for examples of using the helpers.
