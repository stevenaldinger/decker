// environment variable config
variable "target_host" {
  type = "string"
}

// general resource gathering
resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  plugin_enabled = "true"
  // use the local DNS server since no external one will recognize "dvwa"
  dns_server = "127.0.0.11"
}
resource "nmap_simple" "nmap" {
  for_each = "${nslookup.ip_address}"
  host = "${each.key}"
  plugin_enabled = "true"
}

// Web server vulnerability scan
resource "nikto" "nikto" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}

// w3af "all usage" scan
resource "w3af" "w3af" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
