// environment variable config
variable "target_host" {
  type = "string"
}

resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  plugin_enabled = "true"
  dns_server = "8.8.4.4"
}

resource "metasploit" "metasploit" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/portscan/tcp"
  options = {
    RHOSTS = "${each.key}/32"
    INTERFACE = "eth0"
  }
}
