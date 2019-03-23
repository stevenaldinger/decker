input "say_hello_to" {
  type = "string"
  default = "world"
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}

// outputs the plugin will return
output "raw_output" {
  type = "string"
}
