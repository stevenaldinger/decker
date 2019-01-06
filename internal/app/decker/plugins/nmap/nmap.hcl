// inputs must be given in the main config spec
// if no default is given, considered to be required
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
// "example.com"
output "host" {
  type = "string"
}
// "123.456.789.1"
output "host_address" {
  type = "string"
}
// port range "1"-"30000" - "open" or "closed"
output "1-30000" {
  type = "string"
}
