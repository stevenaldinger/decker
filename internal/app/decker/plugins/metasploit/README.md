# metasploit decker plugin

[metasploit](https://www.metasploit.com/)

## Example usage

```hcl
resource "metasploit" "metasploit" {
  for_each = "${nslookup.ip_address}"
  exploit = "auxiliary/scanner/portscan/tcp"
  options = {
    RHOSTS = "${each.key}/32"
    INTERFACE = "eth0"
  }
}
```
