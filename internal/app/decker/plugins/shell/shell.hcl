// inputs must be given in the main config spec
// if no default is given, considered to be required
input "command" {
  type = "string"
  default = ""
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}

input "for_each" {
  type = "list"
  default = []
}
