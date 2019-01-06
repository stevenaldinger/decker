// environment variable config
variable "target_host" {
  type = "string"
}

// general resource gathering
resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  plugin_enabled = "true"
  dns_server = "8.8.4.4"
}
resource "nmap" "nmap" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}

resource "sslscan" "sslscan" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}
resource "sslyze" "sslyze" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}

// WAF detection and bypassing
resource "wafw00f" "wafw00f" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "wafninja" "wafninja" {
  host = "${var.target_host}"
  plugin_enabled = "${wafw00f.waf_detected}"
}
