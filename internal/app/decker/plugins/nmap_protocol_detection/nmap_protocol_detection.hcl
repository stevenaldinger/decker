// inputs must be given in the main config spec
// if no default is given, considered to be required
// input "arguments" {
//   type = "list"
//   default = []
// }
input "type" {
  type = "string"
  default = "port_scan"
}

input "host" {
  type = "string"
  default = "example.com"
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}

// outputs the plugin will return
output "raw_output" {
  type = "string"
}
