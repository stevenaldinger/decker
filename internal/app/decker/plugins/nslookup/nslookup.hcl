// inputs must be given in the main config spec
// if no default is given, considered to be required
input "host" {
  type = "string"
  default = "example.com"
}

input "dns_server" {
  type = "string"
  default = "8.8.8.8"
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}

// "8.8.4.4"
output "dns_server" {
  type = "string"
}
// "8.8.4.4#53"
output "dns_address" {
  type = "string"
}
// "example.com"
output "host_name" {
  type = "string"
}
// "172.217.11.142"
output "ip_address" {
  type = "string"
}
