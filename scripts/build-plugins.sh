#!/bin/bash

plugins_dir=""

if [ -z "$1" ]
then
  echo "You need to pass in an argument to point to the plugins directory"
else
 plugins_dir="$1"
fi

for file in "$plugins_dir"/*
do
  plugin_name="$(basename $file)"

  echo "Building plugin: $plugin_name"

	CGO_ENABLED=1 go build -buildmode=plugin -o "$plugins_dir/$plugin_name/$plugin_name.so" "$plugins_dir/$plugin_name/main.go"
done
