input "exploit" {
  type = "string"
  default = ""
}

input "options" {
  type = "map"
  default = {}
}

// input "db_enabled" {
//   type = "string"
//   default = "false"
// }

input "plugin_enabled" {
  type = "string"
  default = "true"
}

// outputs the plugin will return
output "raw_output" {
  type = "string"
}
