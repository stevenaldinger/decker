// environment variable config
variable "say_hello" {
  type = "string"
}

// says hello to a given string input
resource "hello_world" "my_plugin" {
  say_hello_to = "${var.say_hello}"
}

resource "hello_world" "my_plugin_2" {
  say_hello_to = "${my_plugin.said_hello_to}... again!"
}
