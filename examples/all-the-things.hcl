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
resource "nslookup" "nslookup" {
  host = "${var.target_host}"
  dns_server = "8.8.4.4"
}
resource "nmap" "nmap" {
  host = "${var.target_host}"
}
resource "nmap" "nmap_all" {
  for_each = "${nslookup.ip_address}"
  host = "${each.key}"
}

resource "dig" "dig" {
  host = "${var.target_host}"
}
resource "dnsrecon" "dnsrecon" {
  host = "${var.target_host}"
}
resource "sublist3r" "sublist3r" {
  host = "${var.target_host}"
}
resource "the_harvester" "the_harvester" {
  host = "${var.target_host}"
}

// SSL vulnerability scans
resource "sslyze" "sslyze" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}
resource "sslscan" "sslscan" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}
resource "a2sv" "a2sv" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.443 == "open"}"
}

// XSS vulnerability scans
resource "xss" "xss" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.80 == "open" || nmap.443 == "open"}"
}

// Web server vulnerability scan
resource "nikto" "nikto" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.80 == "open" || nmap.443 == "open"}"
}

// w3af "all usage" scan
resource "w3af" "w3af" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.80 == "open" || nmap.443 == "open"}"
}

// WAF detection and bypassing
resource "wafw00f" "wafw00f" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.80 == "open" || nmap.443 == "open"}"
}
resource "wafninja" "wafninja" {
  host = "${var.target_host}"
  plugin_enabled = "${nmap.80 == "open"  && wafw00f.waf_detected}"
}

// SMTP scans
resource "metasploit" "smtp_enum" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/smtp/smtp_enum"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].25 == "open"}"
}

// FTP scans
resource "metasploit" "ftp_version" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/ftp/ftp_version"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].21 == "open"}"
}
resource "metasploit" "ftp_anonymous_login" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/ftp/anonymous"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].21 == "open"}"
}

// MySQL scans
resource "metasploit" "mysql_version" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/mysql/mysql_version"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].3306 == "open"}"
}

// Postgres scans
resource "metasploit" "postgres_version" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/postgres/postgres_version"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].5432 == "open"}"
}

// IMAP scans
resource "metasploit" "imap_version" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/imap/imap_version"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].143 == "open"}"
}

// VNC scans
resource "metasploit" "vnc_anonymous_login" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/vnc/vnc_none_auth"
  options = {
    RHOSTS = "${each.key}"
  }
  plugin_enabled = "${nmap_all["${each.key}"].5900 == "open"}"
}
