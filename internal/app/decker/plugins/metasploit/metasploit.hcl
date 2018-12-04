input "exploit" {
  type = "string"
  default = ""
}

input "options" {
  type = "map"
  default = {}
}

input "for_each" {
  type = "list"
  default = []
}

input "plugin_enabled" {
  type = "string"
  default = "true"
}
