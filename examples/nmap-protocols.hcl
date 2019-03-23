// environment variable config
variable "target_host" {
  type = "string"
}

resource "nmap_protocol_detection" "nmap" {
  host = "${var.target_host}"
  type = "protocol_detection"
}

resource "metasploit" "metasploit" {
  for_each = "${nmap.ssh}"
  exploit = "auxiliary/scanner/ssh/ssh_login"
  options = {
    RHOSTS = "${var.target_host}"
    RPORT = "${each.key}"
    USERPASS_FILE = "/usr/share/metasploit-framework/data/wordlists/root_userpass.txt"
  }
}

resource "sslscan" "sslscan" {
  for_each = "${nmap.https}"
  host = "${var.target_host}"
}
