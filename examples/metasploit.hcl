// environment variable config
variable "target_host" {
  type = "string"
}

resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  plugin_enabled = "true"
  dns_server = "8.8.4.4"
}

resource "nmap" "nmap" {
  for_each = "${nslookup.ip_address}"
  host = "${each.key}"
  plugin_enabled = "true"
}

// for each IP, check if nmap found port 25 open.
// if yes, run metasploit's smtp_enum scanner
resource "metasploit" "metasploit" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/smtp/smtp_enum"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap["${each.key}"].25 == "open"}"
}
