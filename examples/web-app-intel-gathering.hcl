// variables are pulled from environment
//   ex: DECKER_TARGET_HOST
// they will be available throughout the config files as var.*
//   ex: ${var.target_host}
variable "target_host" {
  type = "string"
}

// resources refer to plugins
// resources need unique names for use in dependency graph
// their outputs will be available to others using the form resource_name.*
//   ex: nslookup_1.ip

// general info gathering
resource "nslookup" "nslookup_1" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "dig" "dig" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "dnsrecon" "dnsrecon" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "nmap" "nmap" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "sublist3r" "sublist3r" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "the_harvester" "the_harvester" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}

// SSL vulnerability scans
resource "sslyze" "sslyze" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "sslscan" "sslscan" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
resource "a2sv" "a2sv" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}

// XSS vulnerability scans
resource "xss" "xss" {
  host = "${var.target_host}"
  plugin_enabled = "true"
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

// WordPress vulnerability scan
resource "wpscan" "wpscan" {
  host = "${var.target_host}"
  plugin_enabled = "true"
}
