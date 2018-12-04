// inputs must be given in the main config spec
// if no default is given, considered to be required
input "host" {
  type = "string"
  default = "example.com"
}

input "for_each" {
  type = "list"
  default = []
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}

// "example.com"
output "host" {
  type = "string"
}
// "123.456.789.1"
output "host_address" {
  type = "string"
}
// "open"
output "22" {
  type = "string"
}
// "open"
output "443" {
  type = "string"
}
