# nslookup decker plugin

[nslookup](https://en.wikipedia.org/wiki/Nslookup)

## Example usage

```
resource "nslookup" "nslookup_1" {
  host = "${var.target_host}"
  plugin_enabled = "true"
  dns_server = "8.8.8.8"
}
```
