# Building Go Plugins

```sh
CGO_ENABLED=1 go build -buildmode=plugin -o plugin.so "$plugins_dir/main.go"`
```
